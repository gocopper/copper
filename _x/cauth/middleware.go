package cauth

import (
	"context"
	"errors"
	"net/http"

	"github.com/tusharsoni/copper/cacl"

	"github.com/tusharsoni/copper/chttp"
	"github.com/tusharsoni/copper/clogger"
	"go.uber.org/fx"
)

type Middleware interface {
	VerifySessionToken(next http.Handler) http.Handler
}

type middleware struct {
	rw     chttp.ReaderWriter
	svc    Svc
	acl    cacl.Svc
	logger clogger.Logger
}

type MiddlewareParams struct {
	fx.In

	RW     chttp.ReaderWriter
	Svc    Svc
	Logger clogger.Logger

	ACL cacl.Svc `optional:"true"`
}

func NewAuthMiddleware(p MiddlewareParams) Middleware {
	return &middleware{
		rw:     p.RW,
		svc:    p.Svc,
		acl:    p.ACL,
		logger: p.Logger,
	}
}

func (m *middleware) VerifySessionToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userUUID, sessionToken, ok := r.BasicAuth()
		if !ok || userUUID == "" || sessionToken == "" {
			m.rw.Unauthorized(w)
			return
		}

		ok, err := m.svc.VerifySessionToken(r.Context(), userUUID, sessionToken)
		if err != nil {
			m.logger.WithTags(map[string]interface{}{
				"userUUID": userUUID,
			}).Error("Failed to verify session token", err)
			m.rw.InternalErr(w)
			return
		}
		if !ok {
			m.rw.Unauthorized(w)
			return
		}

		impersonatedUserUUID := r.Header.Get("x-user-uuid")
		if impersonatedUserUUID != "" {
			if m.acl == nil {
				m.logger.Error("Failed to impersonate user", errors.New("acl is not configured"))
				m.rw.InternalErr(w)
				return
			}

			ok, err := m.acl.UserHasPermission(r.Context(), userUUID, "cauth/session", "impersonate")
			if err != nil {
				m.logger.WithTags(map[string]interface{}{
					"userUUID": userUUID,
				}).Error("Failed to impersonate user", err)
				m.rw.InternalErr(w)
				return
			}

			if !ok {
				m.logger.WithTags(map[string]interface{}{
					"userUUID": userUUID,
				}).Error("Failed to impersonate user", errors.New("user does not have permission to impersonate"))
				m.rw.InternalErr(w)
				return
			}

			userUUID = impersonatedUserUUID
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeySession, userUUID)))
	})
}
