package anonymous

import (
	"context"

	"github.com/tusharsoni/copper/cauth"
	"github.com/tusharsoni/copper/cerror"
	"go.uber.org/fx"
)

type Svc interface {
	CreateAnonymousUser(ctx context.Context) (u *cauth.User, sessionToken string, err error)
}

type SvcParams struct {
	fx.In

	Auth cauth.Svc
}

func NewSvc(p SvcParams) Svc {
	return &svc{
		auth: p.Auth,
	}
}

type svc struct {
	auth cauth.Svc
}

func (s *svc) CreateAnonymousUser(ctx context.Context) (user *cauth.User, sessionToken string, err error) {
	user, err = s.auth.CreateUser(ctx)
	if err != nil {
		return nil, "", cerror.New(err, "failed to create user", nil)
	}

	sessionToken, err = s.auth.ResetSessionToken(ctx, user.UUID)
	if err != nil {
		return nil, "", cerror.New(err, "failed to reset session token", map[string]interface{}{
			"userUUID": user.UUID,
		})
	}

	return user, sessionToken, nil
}
