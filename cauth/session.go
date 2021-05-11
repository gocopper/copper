package cauth

import (
	"context"
	"net/http"

	"github.com/gocopper/copper/chttp"
	"github.com/gocopper/copper/clogger"
)

type ctxKey string

const ctxKeySession = ctxKey("cauth/session")

// NewSetSessionMiddleware instantiates and creates a new SetSessionMiddleware
func NewSetSessionMiddleware(auth *Svc, rw *chttp.ReaderWriter, logger clogger.Logger) *SetSessionMiddleware {
	return &SetSessionMiddleware{
		auth:   auth,
		rw:     rw,
		logger: logger,
	}
}

// SetSessionMiddleware is a middleware that checks for a valid session uuid and token in the Authorization header
// using basic auth. If the session is present, it is validated, saved in the request ctx, and the next handler is
// called. If the session is invalid, an unauthorized response is sent back.
// The next handler is also called if the authorizaiton header is missing. To ensure verified session, use in conjunction
// with VerifySessionMiddleware.
type SetSessionMiddleware struct {
	auth   *Svc
	rw     *chttp.ReaderWriter
	logger clogger.Logger
}

// Handle implements the middleware for SetSessionMiddleware.
func (mw *SetSessionMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionUUID, plainToken, ok := r.BasicAuth()
		if !ok || sessionUUID == "" || plainToken == "" {
			next.ServeHTTP(w, r)

			return
		}

		ok, session, err := mw.auth.ValidateSession(r.Context(), sessionUUID, plainToken)
		if err != nil {
			mw.logger.WithTags(map[string]interface{}{
				"sessionUUID": sessionUUID,
			}).Error("Failed to verify session token", err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		if !ok {
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeySession, session)))
	})
}

// NewVerifySessionMiddleware instantiates and creates a new VerifySessionMiddleware.
func NewVerifySessionMiddleware(auth *Svc, rw *chttp.ReaderWriter, logger clogger.Logger) *VerifySessionMiddleware {
	return &VerifySessionMiddleware{
		auth:   auth,
		rw:     rw,
		logger: logger,
	}
}

// VerifySessionMiddleware is a middleware that checks for a valid session object in the request context. The session
// can be set in the request context with the SetSessionMiddleware.
type VerifySessionMiddleware struct {
	auth   *Svc
	rw     *chttp.ReaderWriter
	logger clogger.Logger
}

// Handle implements the middleware for VerifySessionMiddleware.
func (mw *VerifySessionMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ok := HasVerifiedSession(r.Context())
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		next.ServeHTTP(w, r)
	})
}

// GetCurrentSession returns the session in the HTTP request context. It should only be used in HTTP request
// handlers that have the SetSessionMiddleware on them. If a session is not found, this method will panic. To avoid
// panics, verify that a session exists either with the VerifySessionMiddleware or the HasVerifiedSession function.
func GetCurrentSession(ctx context.Context) *Session {
	session, ok := ctx.Value(ctxKeySession).(*Session)
	if !ok || session == nil {
		panic("session not found in context")
	}

	return session
}

// HasVerifiedSession checks if the context has a valid session
func HasVerifiedSession(ctx context.Context) bool {
	session, ok := ctx.Value(ctxKeySession).(*Session)

	return ok && session != nil
}
