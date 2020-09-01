package cauth

import (
	"context"
	"net/http"

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
	logger clogger.Logger
}

type MiddlewareParams struct {
	fx.In

	RW     chttp.ReaderWriter
	Svc    Svc
	Logger clogger.Logger
}

func NewAuthMiddleware(p MiddlewareParams) Middleware {
	return &middleware{
		rw:     p.RW,
		svc:    p.Svc,
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
		if impersonatedUserUUID != "" { // todo: check for impersonation permissions
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeySession, impersonatedUserUUID)))
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeySession, userUUID)))
	})
}
