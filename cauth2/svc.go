package cauth2

import (
	"context"

	"github.com/google/uuid"
	"github.com/tusharsoni/copper/cerror"
)

type Svc interface {
	CreateUser(ctx context.Context) (*User, error)
	VerifySessionToken(ctx context.Context, userUUID, token string) (bool, error)
	ResetSessionToken(ctx context.Context, userUUID string) (string, error)
}

type svc struct {
	repo Repo
}

func NewSvc(repo Repo) Svc {
	return &svc{
		repo: repo,
	}
}

func (s *svc) CreateUser(ctx context.Context) (*User, error) {
	_, encrypted, err := generateSessionToken()
	if err != nil {
		return nil, cerror.New(err, "failed to generate session token", nil)
	}

	u := User{
		UUID:         uuid.New().String(),
		SessionToken: encrypted,
	}

	err = s.repo.AddUser(ctx, &u)
	if err != nil {
		return nil, cerror.New(err, "failed to insert user", nil)
	}

	return &u, nil
}

func (s *svc) VerifySessionToken(ctx context.Context, userUUID, token string) (bool, error) {
	u, err := s.repo.GetUser(ctx, userUUID)
	if err != nil {
		return false, cerror.New(err, "failed to get user", map[string]interface{}{
			"uuid": userUUID,
		})
	}

	return verifySessionToken(u.SessionToken, token), nil
}

func (s *svc) ResetSessionToken(ctx context.Context, userUUID string) (string, error) {
	u, err := s.repo.GetUser(ctx, userUUID)
	if err != nil {
		return "", cerror.New(err, "failed to get user", map[string]interface{}{
			"uuid": userUUID,
		})
	}

	raw, encrypted, err := generateSessionToken()
	if err != nil {
		return "", cerror.New(err, "failed to generate session token", nil)
	}

	u.SessionToken = encrypted

	err = s.repo.AddUser(ctx, u)
	if err != nil {
		return "", cerror.New(err, "failed to update user's session token", map[string]interface{}{
			"userId": u.UUID,
		})
	}

	return raw, nil
}
