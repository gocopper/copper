package emailotp

import (
	"context"
	"html/template"
	"strings"

	"github.com/tusharsoni/copper/cmailer"

	"github.com/tusharsoni/copper/cauth"

	"github.com/jinzhu/gorm"
	"github.com/tusharsoni/copper/cerror"
	"github.com/tusharsoni/copper/crandom"
)

type Svc interface {
	Signup(ctx context.Context, email string) (c *Credentials, token *string, err error)
	Login(ctx context.Context, email string, verificationCode uint) (c *Credentials, token string, err error)
}

type svc struct {
	auth   cauth.Svc
	repo   Repo
	mailer cmailer.Mailer
	config Config

	verificationEmailTemplate *template.Template
}

func NewSvc(auth cauth.Svc, repo Repo, mailer cmailer.Mailer, config Config) (Svc, error) {
	verificationEmailTemplate, err := template.New("VerificationEmail").Parse(config.VerificationEmail.BodyTemplate)
	if err != nil {
		return nil, cerror.New(err, "failed to parse verification email template", nil)
	}

	return &svc{
		auth:   auth,
		repo:   repo,
		mailer: mailer,
		config: config,

		verificationEmailTemplate: verificationEmailTemplate,
	}, nil
}

func (s *svc) Signup(ctx context.Context, email string) (*Credentials, *string, error) {
	var (
		newUser               = false
		verificationEmailBody strings.Builder
	)

	c, err := s.repo.GetCredentialsByEmail(ctx, email)
	if err != nil && !cerror.HasCause(err, gorm.ErrRecordNotFound) {
		return nil, nil, cerror.New(err, "failed to get credentials", map[string]interface{}{
			"email": email,
		})
	}

	if c == nil {
		newUser = true

		u, err := s.auth.CreateUser(ctx)
		if err != nil {
			return nil, nil, cerror.New(err, "failed to create new user", nil)
		}

		c = &Credentials{
			UserUUID: u.UUID,
			Email:    email,
		}
	}

	c.VerificationCode = uint(crandom.GenerateRandomNumericalCode(4))
	c.Verified = false

	err = s.repo.AddCredentials(ctx, c)
	if err != nil {
		return nil, nil, cerror.New(err, "failed to update credentials", map[string]interface{}{
			"userUUID": c.UserUUID,
		})
	}

	err = s.verificationEmailTemplate.Execute(&verificationEmailBody, map[string]interface{}{
		"VerificationCode": c.VerificationCode,
	})
	if err != nil {
		return nil, nil, cerror.New(err, "failed to execute verification email template", nil)
	}

	_, err = s.mailer.SendHTML(ctx,
		s.config.VerificationEmail.From,
		c.Email,
		s.config.VerificationEmail.Subject,
		verificationEmailBody.String(),
	)
	if err != nil {
		return nil, nil, cerror.New(err, "failed to email verification code", map[string]interface{}{
			"userUUID": c.UserUUID,
		})
	}

	if !newUser || s.config.RequiresVerification {
		return c, nil, nil
	}

	sessionToken, err := s.auth.ResetSessionToken(ctx, c.UserUUID)
	if err != nil {
		return nil, nil, cerror.New(err, "failed to reset session token", map[string]interface{}{
			"userUUID": c.UserUUID,
		})
	}

	return c, &sessionToken, nil
}

func (s *svc) Login(ctx context.Context, email string, verificationCode uint) (*Credentials, string, error) {
	c, err := s.repo.GetCredentialsByEmail(ctx, email)
	if err != nil {
		return nil, "", cerror.New(err, "failed to get credentials", map[string]interface{}{
			"email": email,
		})
	}

	if c.VerificationCode != verificationCode {
		return nil, "", cerror.New(nil, "invalid verification code", nil)
	}

	c.Verified = true

	err = s.repo.AddCredentials(ctx, c)
	if err != nil {
		return nil, "", cerror.New(err, "failed to update credentials", map[string]interface{}{
			"userUUID": c.UserUUID,
		})
	}

	token, err := s.auth.ResetSessionToken(ctx, c.UserUUID)
	if err != nil {
		return nil, "", cerror.New(err, "failed to reset session token", map[string]interface{}{
			"userUUID": c.UserUUID,
		})
	}

	return c, token, nil
}
