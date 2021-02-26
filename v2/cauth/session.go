package cauth

import (
	"context"
	"net/http"

	"github.com/tusharsoni/copper/v2/chttp"
	"github.com/tusharsoni/copper/v2/clogger"
)

type ctxKey string

const ctxKeySession = ctxKey("cauth/session")

// NewVerifySessionMiddleware instantiates and creates a new VerifySessionMiddleware.
func NewVerifySessionMiddleware(auth *Svc, rw chttp.ReaderWriter, logger clogger.Logger) *VerifySessionMiddleware {
	return &VerifySessionMiddleware{
		auth:   auth,
		rw:     rw,
		logger: logger,
	}
}

// VerifySessionMiddleware is a middleware that checks for a valid session uuid and token in the Authorization header
// using basic auth. If the session is valid, the session object is saved in the request ctx and the next handler is
// called. If the session is invalid, an unauthorized response is sent back.
type VerifySessionMiddleware struct {
	auth   *Svc
	rw     chttp.ReaderWriter
	logger clogger.Logger
}

// Handle implements the middleware for VerifySessionMiddleware.
func (mw *VerifySessionMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionUUID, plainToken, ok := r.BasicAuth()
		if !ok || sessionUUID == "" || plainToken == "" {
			mw.rw.Unauthorized(w)

			return
		}

		ok, session, err := mw.auth.ValidateSession(r.Context(), sessionUUID, plainToken)
		if err != nil {
			mw.logger.WithTags(map[string]interface{}{
				"sessionUUID": sessionUUID,
			}).Error("Failed to verify session token", err)
			mw.rw.InternalErr(w)

			return
		}

		if !ok {
			mw.rw.Unauthorized(w)

			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeySession, session)))
	})
}

// GetCurrentSession returns the session in the HTTP request context. It should only be used in HTTP request
// handlers that have the VerifySessionMiddleware on them. If a session is not found, this method will panic.
func GetCurrentSession(ctx context.Context) *Session {
	session, ok := ctx.Value(ctxKeySession).(*Session)
	if !ok || session == nil {
		panic("session not found in context")
	}

	return session
}
