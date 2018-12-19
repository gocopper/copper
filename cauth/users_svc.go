package cauth

import (
	"context"
	"errors"
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

// UsersSvc provides high level methods to manage users.
type UsersSvc interface {
	Signup(ctx context.Context, email, password string) (*user, error)
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

func (s *usersSvc) Signup(ctx context.Context, email, password string) (*user, error) {
	_, err := s.users.FindByEmail(ctx, email)
	if err != nil && cerror.Cause(err) != ErrUserNotFound {
		return nil, cerror.New(err, "failed to find u by email", nil)
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
