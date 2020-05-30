package email

import (
	"context"
	"html/template"
	"strings"

	"github.com/tusharsoni/copper/cauth"
	"github.com/tusharsoni/copper/cerror"
	"github.com/tusharsoni/copper/clogger"
	"github.com/tusharsoni/copper/cmailer"
	"github.com/tusharsoni/copper/crandom"
	"go.uber.org/fx"
	"golang.org/x/crypto/bcrypt"
)

type Svc interface {
	Signup(ctx context.Context, email, password string) (c *Credentials, sessionToken string, err error)
	Login(ctx context.Context, email, password string) (c *Credentials, sessionToken string, err error)
	VerifyUser(ctx context.Context, uuid string, verificationCode string) error
	ResendVerificationCode(ctx context.Context, uuid string) error
	ResetPassword(ctx context.Context, email string) error
	ChangePassword(ctx context.Context, email, oldPassword, newPassword string) error
}

type svc struct {
	auth   cauth.Svc
	repo   Repo
	mailer cmailer.Mailer
	config Config
	logger clogger.Logger
}

type SvcParams struct {
	fx.In

	Auth   cauth.Svc
	Repo   Repo
	Mailer cmailer.Mailer
	Config Config
	Logger clogger.Logger
}

func NewSvc(p SvcParams) Svc {
	return &svc{
		auth:   p.Auth,
		repo:   p.Repo,
		mailer: p.Mailer,
		config: p.Config,
		logger: p.Logger,
	}
}

func (s *svc) Signup(ctx context.Context, email, password string) (c *Credentials, sessionToken string, err error) {
	_, err = s.repo.GetCredentialsByEmail(ctx, email)
	if err != nil && cerror.Cause(err) != cauth.ErrUserNotFound {
		return nil, "", cerror.New(err, "failed to find credentials", map[string]interface{}{
			"email": email,
		})
	}
	if err == nil {
		return nil, "", cauth.ErrUserAlreadyExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), s.config.PasswordHashCost)
	if err != nil {
		return nil, "", cerror.New(err, "failed to generate password hash", map[string]interface{}{
			"email": email,
		})
	}

	u, err := s.auth.CreateUser(ctx)
	if err != nil {
		return nil, "", cerror.New(err, "failed to create user", nil)
	}

	c = &Credentials{
		UserUUID:         u.UUID,
		Email:            email,
		Password:         string(passwordHash),
		Verified:         false,
		VerificationCode: crandom.GenerateRandomString(s.config.VerificationCodeLen),
	}

	err = s.repo.AddCredentials(ctx, c)
	if err != nil {
		return nil, "", cerror.New(err, "failed to insert credentials", map[string]interface{}{
			"email":            c.Email,
			"verificationCode": c.VerificationCode,
		})
	}

	sessionToken, err = s.auth.ResetSessionToken(ctx, c.UserUUID)
	if err != nil {
		return nil, "", cerror.New(err, "failed to reset session token", map[string]interface{}{
			"userUUID": c.UserUUID,
		})
	}

	go func() {
		err = s.sendVerificationCodeEmail(c)
		if err != nil {
			s.logger.WithTags(map[string]interface{}{
				"userUUID": c.UserUUID,
			}).Error("Failed to send verification code email", err)
		}
	}()

	return c, sessionToken, nil
}

func (s *svc) Login(ctx context.Context, email, password string) (c *Credentials, sessionToken string, err error) {
	c, err = s.repo.GetCredentialsByEmail(ctx, email)
	if err != nil && cerror.Cause(err) != cauth.ErrUserNotFound {
		return nil, "", cerror.New(err, "failed to find credentials", map[string]interface{}{
			"email": email,
		})
	} else if cerror.Cause(err) == cauth.ErrUserNotFound {
		return nil, "", cauth.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(c.Password), []byte(password))
	if err != nil {
		return nil, "", cauth.ErrInvalidCredentials
	}

	sessionToken, err = s.auth.ResetSessionToken(ctx, c.UserUUID)
	if err != nil {
		return nil, "", cerror.New(err, "failed to reset session token", map[string]interface{}{
			"userUUID": c.UserUUID,
		})
	}

	return c, sessionToken, nil
}

func (s *svc) VerifyUser(ctx context.Context, uuid string, verificationCode string) error {
	c, err := s.repo.GetCredentialsByUserUUID(ctx, uuid)
	if err != nil {
		return cerror.New(err, "failed to find user", map[string]interface{}{
			"uuid": uuid,
		})
	}

	if c.VerificationCode != verificationCode {
		return cauth.ErrInvalidCredentials
	}

	c.Verified = true

	err = s.repo.AddCredentials(ctx, c)
	if err != nil {
		return cerror.New(err, "failed to save credentials", map[string]interface{}{
			"userUUID": c.UserUUID,
		})
	}

	return nil
}

func (s *svc) ResendVerificationCode(ctx context.Context, uuid string) error {
	c, err := s.repo.GetCredentialsByUserUUID(ctx, uuid)
	if err != nil {
		return cerror.New(err, "failed to get credentials", map[string]interface{}{
			"userUUID": uuid,
		})
	}

	err = s.sendVerificationCodeEmail(c)
	if err != nil {
		return cerror.New(err, "failed to send verification code email", map[string]interface{}{
			"userUUID": uuid,
		})
	}

	return nil
}

func (s *svc) ResetPassword(ctx context.Context, email string) error {
	c, err := s.repo.GetCredentialsByEmail(ctx, email)
	if err != nil {
		return cerror.New(err, "failed to find credentials", map[string]interface{}{
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

	c.Password = string(resetPasswordTokenHash)

	err = s.repo.AddCredentials(ctx, c)
	if err != nil {
		return cerror.New(err, "failed to update credentials with reset password token", map[string]interface{}{
			"userUUID": c.UserUUID,
		})
	}

	go func() {
		err = s.sendResetPasswordTokenEmail(c, resetPasswordToken)
		if err != nil {
			s.logger.WithTags(map[string]interface{}{
				"userUUID": c.UserUUID,
			}).Error("Failed to send reset password token email", err)
		}
	}()

	return nil
}

func (s *svc) ChangePassword(ctx context.Context, email, oldPassword, newPassword string) error {
	c, err := s.repo.GetCredentialsByEmail(ctx, email)
	if err != nil && cerror.Cause(err) != cauth.ErrUserNotFound {
		return cerror.New(err, "failed to find credentials", map[string]interface{}{
			"email": email,
		})
	} else if cerror.Cause(err) == cauth.ErrUserNotFound {
		return cauth.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(c.Password), []byte(oldPassword))
	if err != nil {
		return cauth.ErrInvalidCredentials
	}

	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), s.config.PasswordHashCost)
	if err != nil {
		return cerror.New(err, "failed to generate hash for new password", nil)
	}

	c.Password = string(newPasswordHash)

	err = s.repo.AddCredentials(ctx, c)
	if err != nil {
		return cerror.New(err, "failed to update credentials with new password", nil)
	}

	return nil
}

func (s *svc) sendResetPasswordTokenEmail(u *Credentials, resetToken string) error {
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

	_, err = s.mailer.SendPlain(
		s.config.ResetPasswordEmail.From,
		u.Email,
		s.config.ResetPasswordEmail.Subject,
		body.String(),
	)
	if err != nil {
		return cerror.New(err, "failed to send reset password email", nil)
	}

	return nil
}

func (s *svc) sendVerificationCodeEmail(u *Credentials) error {
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

	_, err = s.mailer.SendPlain(
		s.config.VerificationEmail.From,
		u.Email,
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
