package cauth

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"text/template"

	"github.com/tusharsoni/copper/clogger"

	"github.com/tusharsoni/copper/cmailer"

	"github.com/tusharsoni/copper/crandom"

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
	Login(ctx context.Context, email, password string) (user *user, sessionToken string, err error)
	Signup(ctx context.Context, email, password string) (*user, error)
	VerifySessionToken(ctx context.Context, email, token string) (*user, error)
	VerifyUser(ctx context.Context, userID uint, verificationCode string) error
	ResendVerificationCode(ctx context.Context, userID uint) error
}

type usersSvc struct {
	users  UserRepo
	mailer cmailer.Mailer
	config Config
	logger clogger.Logger
}

func newUsersSvc(users UserRepo, mailer cmailer.Mailer, config Config, logger clogger.Logger) UsersSvc {
	return &usersSvc{
		users:  users,
		mailer: mailer,
		config: config,
		logger: logger,
	}
}

func (s *usersSvc) VerifyUser(ctx context.Context, userID uint, verificationCode string) error {
	u, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return cerror.New(err, "failed to find user by id", map[string]string{
			"id": strconv.Itoa(int(userID)),
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

	return nil
}

func (s *usersSvc) VerifySessionToken(ctx context.Context, email, token string) (*user, error) {
	u, err := s.users.FindByEmail(ctx, email)
	if err != nil && cerror.Cause(err) != ErrUserNotFound {
		return nil, cerror.New(err, "failed to find user by email", nil)
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
	} else if err == ErrUserNotFound {
		return nil, "", ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	sessionToken, err = s.resetSessionToken(ctx, u)
	if err != nil {
		return nil, "", cerror.New(err, "failed to reset user's session token", nil)
	}

	return u, sessionToken, nil
}

func (s *usersSvc) ResendVerificationCode(ctx context.Context, userID uint) error {
	u, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return cerror.New(err, "failed to get user by id", nil)
	}

	err = s.sendVerificationCodeEmail(u)
	if err != nil {
		return cerror.New(err, "failed to send verification code email", nil)
	}

	return nil
}

func (s *usersSvc) Signup(ctx context.Context, email, password string) (*user, error) {
	_, err := s.users.FindByEmail(ctx, email)
	if err != nil && cerror.Cause(err) != ErrUserNotFound {
		return nil, cerror.New(err, "failed to find user by email", nil)
	}
	if err == nil {
		return nil, ErrUserAlreadyExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), s.config.PasswordHashCost)
	if err != nil {
		return nil, cerror.New(err, "failed to generate password hash", map[string]string{
			"email": email,
		})
	}

	u := user{
		Email:            email,
		Password:         string(passwordHash),
		VerificationCode: crandom.GenerateRandomString(s.config.VerificationCodeLen),
	}

	err = s.users.Add(ctx, &u)
	if err != nil {
		return nil, cerror.New(err, "failed to create new user", map[string]string{
			"email":            u.Email,
			"verificationCode": u.VerificationCode,
		})
	}

	go func() {
		err = s.sendVerificationCodeEmail(&u)
		if err != nil {
			s.logger.Error("Failed to send verification code email", err)
		}
	}()

	return &u, nil
}

func (s *usersSvc) sendVerificationCodeEmail(u *user) error {
	var body strings.Builder

	t, err := template.New("cauth/verificationEmail").Parse(s.config.VerificationEmail.BodyTemplate)
	if err != nil {
		return cerror.New(err, "failed to create verification email template", map[string]string{
			"template": s.config.VerificationEmail.BodyTemplate,
		})
	}

	verificationEmailVars := struct {
		VerificationCode string
	}{VerificationCode: u.VerificationCode}

	err = t.Execute(&body, &verificationEmailVars)
	if err != nil {
		return cerror.New(err, "failed to create verification email body", map[string]string{
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
		return cerror.New(err, "failed to send verification email", map[string]string{
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
		return "", cerror.New(err, "failed to generate hashed session token", map[string]string{
			"sessionToken": newToken,
		})
	}

	hashedToken := string(hashedTokenData)
	u.SessionToken = &hashedToken

	err = s.users.Add(ctx, u)
	if err != nil {
		return "", cerror.New(err, "failed to update user's session token", map[string]string{
			"userId":       strconv.Itoa(int(u.ID)),
			"sessionToken": newToken,
		})
	}

	return newToken, nil
}
