package phone

import (
	"context"
	"fmt"

	"github.com/tusharsoni/copper/cauth"

	"github.com/tusharsoni/copper/cerror"
	"github.com/tusharsoni/copper/crandom"
	"github.com/tusharsoni/copper/ctexter"
	"gorm.io/gorm"
)

type Svc interface {
	Signup(ctx context.Context, phoneNumber string) error
	Login(ctx context.Context, phoneNumber string, verificationCode uint) (c *Credentials, token string, err error)
}

type svc struct {
	auth   cauth.Svc
	repo   Repo
	texter ctexter.Svc
}

func NewSvc(auth cauth.Svc, repo Repo, texter ctexter.Svc) Svc {
	return &svc{
		auth:   auth,
		repo:   repo,
		texter: texter,
	}
}

func (s *svc) Signup(ctx context.Context, phoneNumber string) error {
	c, err := s.repo.GetCredentialsByPhoneNumber(ctx, phoneNumber)
	if err != nil && !cerror.HasCause(err, gorm.ErrRecordNotFound) {
		return cerror.New(err, "failed to get credentials", map[string]interface{}{
			"phoneNumber": phoneNumber,
		})
	}

	if c == nil {
		u, err := s.auth.CreateUser(ctx)
		if err != nil {
			return cerror.New(err, "failed to create new user", nil)
		}

		c = &Credentials{
			UserUUID:    u.UUID,
			PhoneNumber: phoneNumber,
		}
	} else {
		_, err = s.auth.ResetSessionToken(ctx, c.UserUUID)
		if err != nil {
			return cerror.New(err, "failed to reset session token", map[string]interface{}{
				"userUUID": c.UserUUID,
			})
		}
	}

	c.VerificationCode = uint(crandom.GenerateRandomNumericalCode(4))
	c.Verified = false

	err = s.repo.AddCredentials(ctx, c)
	if err != nil {
		return cerror.New(err, "failed to update credentials", map[string]interface{}{
			"userUUID": c.UserUUID,
		})
	}

	_, err = s.texter.SendSMS(phoneNumber, fmt.Sprintf("Verification Code: %d", c.VerificationCode))
	if err != nil {
		return cerror.New(err, "failed to text verification code", map[string]interface{}{
			"userUUID": c.UserUUID,
		})
	}

	return nil
}

func (s *svc) Login(ctx context.Context, phoneNumber string, verificationCode uint) (*Credentials, string, error) {
	c, err := s.repo.GetCredentialsByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return nil, "", cerror.New(err, "failed to get credentials", map[string]interface{}{
			"phoneNumber": phoneNumber,
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
