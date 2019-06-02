package cauth

import (
	"context"
	"errors"
	"strings"
	"text/template"
	"time"

	"go.uber.org/fx"

	"github.com/tusharsoni/copper/cpubsub"

	"github.com/tusharsoni/copper/clogger"

	"github.com/tusharsoni/copper/cmailer"

	"github.com/tusharsoni/copper/crandom"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/tusharsoni/copper/cerror"
)

// ErrUserAlreadyExists is returned by UsersSvc when a user already exists. For example, signing up with an email
// with which a user already exists.
var ErrUserAlreadyExists = errors.New("user already exists")

// ErrInvalidCredentials is returned by UsersSvc when the given credentials such as email/password or session token
// are incorrect.
var ErrInvalidCredentials = errors.New("invalid credentials")

// UsersSvc provides high level methods to manage users.
type UsersSvc interface {
	GetByUUID(ctx context.Context, uuid string) (*user, error)
	FindByEmail(ctx context.Context, email string) (*user, error)

	Login(ctx context.Context, email, password string) (user *user, sessionToken string, err error)
	Logout(ctx context.Context, uuid string) error
	Signup(ctx context.Context, email, password string) (user *user, sessionToken string, err error)
	VerifySessionToken(ctx context.Context, uuid, token string) (*user, error)
	VerifyUser(ctx context.Context, uuid string, verificationCode string) error
	ResendVerificationCode(ctx context.Context, uuid string) error
	ResetPassword(ctx context.Context, email string) error
	ChangePassword(ctx context.Context, email, oldPassword, newPassword string) error
}

type usersSvc struct {
	users  userRepo
	mailer cmailer.Mailer
	pubsub cpubsub.PubSub
	config Config
	logger clogger.Logger
}

type usersSvcParams struct {
	fx.In

	Users  userRepo
	Mailer cmailer.Mailer
	PubSub *cpubsub.LocalPubSub
	Config Config
	Logger clogger.Logger
}

func newUsersSvc(p usersSvcParams) UsersSvc {
	return &usersSvc{
		users:  p.Users,
		mailer: p.Mailer,
		pubsub: p.PubSub,
		config: p.Config,
		logger: p.Logger,
	}
}

func (s *usersSvc) GetByUUID(ctx context.Context, uuid string) (*user, error) {
	return s.users.GetByUUID(ctx, uuid)
}

func (s *usersSvc) FindByEmail(ctx context.Context, email string) (*user, error) {
	return s.users.FindByEmail(ctx, email)
}

func (s *usersSvc) ChangePassword(ctx context.Context, email, oldPassword, newPassword string) error {
	u, err := s.users.FindByEmail(ctx, email)
	if err != nil && cerror.Cause(err) != ErrUserNotFound {
		return cerror.New(err, "failed to find user by email", map[string]interface{}{
			"email": email,
		})
	} else if cerror.Cause(err) == ErrUserNotFound {
		return ErrInvalidCredentials
	}

	if u.Password == nil {
		return cerror.New(nil, "user has no password set", map[string]interface{}{
			"email": email,
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(*u.Password), []byte(oldPassword))
	if err != nil {
		return ErrInvalidCredentials
	}

	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), s.config.PasswordHashCost)
	if err != nil {
		return cerror.New(err, "failed to generate hash for new password", nil)
	}

	newPasswordHashStr := string(newPasswordHash)
	u.Password = &newPasswordHashStr

	err = s.users.Add(ctx, u)
	if err != nil {
		return cerror.New(err, "failed to update user with new password", nil)
	}

	return nil
}

func (s *usersSvc) ResetPassword(ctx context.Context, email string) error {
	u, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return cerror.New(err, "failed to find user by email", map[string]interface{}{
			"email": email,
		})
	}

	resetPasswordToken := crandom.GenerateRandomString(s.config.ResetPasswordTokenLen)

	resetPasswordTokenHash, err := bcrypt.GenerateFromPassword([]byte(resetPasswordToken), s.config.PasswordHashCost)
	if err != nil {
		return cerror.New(err, "failed to generate reset password token hash", map[string]interface{}{
			"token": resetPasswordToken,
		})
	}

	resetPasswordTokenHashStr := string(resetPasswordTokenHash)
	u.Password = &resetPasswordTokenHashStr

	err = s.users.Add(ctx, u)
	if err != nil {
		return cerror.New(err, "failed to update user with reset password token", nil)
	}

	go func() {
		err = s.sendResetPasswordTokenEmail(u, resetPasswordToken)
		if err != nil {
			s.logger.Error("Failed to send reset password token email", err)
		}
	}()

	return nil
}

func (s *usersSvc) VerifyUser(ctx context.Context, uuid string, verificationCode string) error {
	u, err := s.users.GetByUUID(ctx, uuid)
	if err != nil {
		return cerror.New(err, "failed to find user by uuid", map[string]interface{}{
			"uuid": uuid,
		})
	}

	if u.VerificationCode != verificationCode {
		return ErrInvalidCredentials
	}

	u.Verified = true

	err = s.users.Add(ctx, u)
	if err != nil {
		return cerror.New(err, "failed to verify user", nil)
	}

	go func() {
		err = s.pubsub.Publish(UserVerifyTopic, []byte(u.UUID))
		if err != nil {
			s.logger.Warn("Failed to publish user verify success update", err)
		}
	}()

	return nil
}

func (s *usersSvc) VerifySessionToken(ctx context.Context, uuid, token string) (*user, error) {
	u, err := s.users.GetByUUID(ctx, uuid)
	if err != nil && cerror.Cause(err) != ErrUserNotFound {
		return nil, cerror.New(err, "failed to find user by uuid", map[string]interface{}{
			"uuid": uuid,
		})
	} else if err == ErrUserNotFound {
		return nil, ErrInvalidCredentials
	}

	if u.SessionToken == nil {
		return nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(*u.SessionToken), []byte(token))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return u, nil
}

func (s *usersSvc) Login(ctx context.Context, email, password string) (user *user, sessionToken string, err error) {
	u, err := s.users.FindByEmail(ctx, email)
	if err != nil && cerror.Cause(err) != ErrUserNotFound {
		return nil, "", cerror.New(err, "failed to find user by email", nil)
	} else if cerror.Cause(err) == ErrUserNotFound {
		return nil, "", ErrInvalidCredentials
	}

	if u.Password == nil {
		return nil, "", cerror.New(nil, "user has no password set", map[string]interface{}{
			"email": email,
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(*u.Password), []byte(password))
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	sessionToken, err = s.resetSessionToken(ctx, u)
	if err != nil {
		return nil, "", cerror.New(err, "failed to reset user's session token", nil)
	}

	now := time.Now()
	u.LastLoginAt = &now

	err = s.users.Add(ctx, u)
	if err != nil {
		s.logger.Warn("Failed to update user's last login at", cerror.WithTags(err, map[string]interface{}{
			"userID": u.UUID,
		}))
	}

	return u, sessionToken, nil
}

func (s *usersSvc) Logout(ctx context.Context, uuid string) error {
	u, err := s.users.GetByUUID(ctx, uuid)
	if err != nil {
		return cerror.New(err, "failed to get user by id", map[string]interface{}{
			"uuid": uuid,
		})
	}

	u.SessionToken = nil

	err = s.users.Add(ctx, u)
	if err != nil {
		return cerror.New(err, "failed to upsert usert", map[string]interface{}{
			"uuid": uuid,
		})
	}

	return nil
}

func (s *usersSvc) ResendVerificationCode(ctx context.Context, uuid string) error {
	u, err := s.users.GetByUUID(ctx, uuid)
	if err != nil {
		return cerror.New(err, "failed to get user by id", nil)
	}

	err = s.sendVerificationCodeEmail(u)
	if err != nil {
		return cerror.New(err, "failed to send verification code email", nil)
	}

	return nil
}

func (s *usersSvc) Signup(ctx context.Context, email, password string) (u *user, sessionToken string, err error) {
	_, err = s.users.FindByEmail(ctx, email)
	if err != nil && cerror.Cause(err) != ErrUserNotFound {
		return nil, "", cerror.New(err, "failed to find user by email", nil)
	}
	if err == nil {
		return nil, "", ErrUserAlreadyExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), s.config.PasswordHashCost)
	if err != nil {
		return nil, "", cerror.New(err, "failed to generate password hash", map[string]interface{}{
			"email": email,
		})
	}
	passwordHashStr := string(passwordHash)

	userUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, "", cerror.New(err, "failed to generate random uuid", nil)
	}

	u = &user{
		UUID:             userUUID.String(),
		Email:            &email,
		Password:         &passwordHashStr,
		VerificationCode: crandom.GenerateRandomString(s.config.VerificationCodeLen),
	}

	err = s.users.Add(ctx, u)
	if err != nil {
		return nil, "", cerror.New(err, "failed to create new user", map[string]interface{}{
			"email":            u.Email,
			"verificationCode": u.VerificationCode,
		})
	}

	go func() {
		err = s.sendVerificationCodeEmail(u)
		if err != nil {
			s.logger.Error("Failed to send verification code email", err)
		}
	}()

	sessionToken, err = s.resetSessionToken(ctx, u)
	if err != nil {
		return nil, "", cerror.New(err, "failed to reset user's session token", map[string]interface{}{
			"userUUID": u.UUID,
		})
	}

	go func() {
		err = s.pubsub.Publish(UserSignupTopic, []byte(u.UUID))
		if err != nil {
			s.logger.Warn("Failed to publish user signup success update", err)
		}
	}()

	return u, sessionToken, nil
}

func (s *usersSvc) sendResetPasswordTokenEmail(u *user, resetToken string) error {
	var body strings.Builder

	t, err := template.New("cauth/resetPasswordTokenEmail").Parse(s.config.ResetPasswordEmail.BodyTemplate)
	if err != nil {
		return cerror.New(err, "failed to create reset password email template", map[string]interface{}{
			"template": s.config.ResetPasswordEmail.BodyTemplate,
		})
	}

	resetPasswordEmailVars := struct {
		ResetToken string
	}{ResetToken: resetToken}

	err = t.Execute(&body, &resetPasswordEmailVars)
	if err != nil {
		return cerror.New(err, "failed to create reset password email body", map[string]interface{}{
			"template": s.config.ResetPasswordEmail.BodyTemplate,
		})
	}

	if u.Email == nil {
		return cerror.New(nil, "user does not have an email", map[string]interface{}{
			"uuid": u.UUID,
		})
	}

	_, err = s.mailer.SendPlain(
		s.config.ResetPasswordEmail.From,
		*u.Email,
		s.config.ResetPasswordEmail.Subject,
		body.String(),
	)
	if err != nil {
		return cerror.New(err, "failed to send reset password email", nil)
	}

	return nil
}

func (s *usersSvc) sendVerificationCodeEmail(u *user) error {
	var body strings.Builder

	t, err := template.New("cauth/verificationEmail").Parse(s.config.VerificationEmail.BodyTemplate)
	if err != nil {
		return cerror.New(err, "failed to create verification email template", map[string]interface{}{
			"template": s.config.VerificationEmail.BodyTemplate,
		})
	}

	verificationEmailVars := struct {
		VerificationCode string
	}{VerificationCode: u.VerificationCode}

	err = t.Execute(&body, &verificationEmailVars)
	if err != nil {
		return cerror.New(err, "failed to create verification email body", map[string]interface{}{
			"template":         s.config.VerificationEmail.BodyTemplate,
			"verificationCode": u.VerificationCode,
		})
	}

	if u.Email == nil {
		return cerror.New(nil, "user does not have an email", map[string]interface{}{
			"uuid": u.UUID,
		})
	}

	_, err = s.mailer.SendPlain(
		s.config.VerificationEmail.From,
		*u.Email,
		s.config.VerificationEmail.Subject,
		body.String(),
	)
	if err != nil {
		return cerror.New(err, "failed to send verification email", map[string]interface{}{
			"from":    s.config.VerificationEmail.From,
			"to":      u.Email,
			"subject": s.config.VerificationEmail.Subject,
			"body":    body.String(),
		})
	}

	return nil
}

func (s *usersSvc) resetSessionToken(ctx context.Context, u *user) (string, error) {
	newToken := crandom.GenerateRandomString(s.config.SessionTokenLen)

	hashedTokenData, err := bcrypt.GenerateFromPassword([]byte(newToken), s.config.PasswordHashCost)
	if err != nil {
		return "", cerror.New(err, "failed to generate hashed session token", map[string]interface{}{
			"sessionToken": newToken,
		})
	}

	hashedToken := string(hashedTokenData)
	u.SessionToken = &hashedToken

	err = s.users.Add(ctx, u)
	if err != nil {
		return "", cerror.New(err, "failed to update user's session token", map[string]interface{}{
			"userId":       u.UUID,
			"sessionToken": newToken,
		})
	}

	return newToken, nil
}
