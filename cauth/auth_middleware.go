package cauth

import (
	"encoding/base64"
	"net/http"
	"strings"

	"go.uber.org/fx"

	"github.com/tusharsoni/copper/chttp"
	"github.com/tusharsoni/copper/clogger"
)

// AuthMiddleware provides a middleware that verifies the auth header using basic auth.
// The username is expected to be the email and the password should be the session token.
// On success, the user is stored in the context.
type AuthMiddleware interface {
	AllowVerified(next http.Handler) http.Handler
	AllowUnverified(next http.Handler) http.Handler
}

type authMiddleware struct {
	resp   chttp.Responder
	users  UsersSvc
	config Config
	logger clogger.Logger
}

type authMiddlewareParams struct {
	fx.In

	Resp   chttp.Responder
	Users  UsersSvc
	Config Config
	Logger clogger.Logger
}

func newAuthMiddleware(p authMiddlewareParams) AuthMiddleware {
	return &authMiddleware{
		resp:   p.Resp,
		users:  p.Users,
		config: p.Config,
		logger: p.Logger,
	}
}

func (m *authMiddleware) AllowVerified(next http.Handler) http.Handler {
	return m.verifyAuth(next, false)
}

func (m *authMiddleware) AllowUnverified(next http.Handler) http.Handler {
	return m.verifyAuth(next, true)
}

func (m *authMiddleware) verifyAuth(next http.Handler, allowUnverified bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email, sessionToken, ok := m.getAuthCredentials(r)
		if !ok {
			m.resp.Unauthorized(w)
			return
		}

		ctx := r.Context()

		user, err := m.users.VerifySessionToken(ctx, email, sessionToken)
		if err != nil && err != ErrInvalidCredentials {
			m.logger.Error("Failed to verify user session token", err)
			m.resp.InternalErr(w)
			return
		} else if err == ErrInvalidCredentials {
			m.resp.Unauthorized(w)
			return
		}

		if !allowUnverified && user.Verified == false {
			m.resp.Unauthorized(w)
			return
		}

		// todo: instead of 'config.AdminEmail', use cacl and check for impersonation permission
		impersonateEmail := r.Header.Get("x-auth-email")
		if impersonateEmail == "" || user.Email != m.config.AdminEmail {
			next.ServeHTTP(w, r.WithContext(ctxWithUser(ctx, user)))
			return
		}

		user, err = m.users.FindByEmail(ctx, impersonateEmail)
		if err != nil {
			m.logger.Error("Failed to find impersonation user by email", err)
			m.resp.InternalErr(w)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctxWithUser(ctx, user)))
	})
}

func (m *authMiddleware) getAuthCredentials(r *http.Request) (username, password string, ok bool) {
	username, password, ok = r.BasicAuth()
	if ok {
		return
	}

	cookie, err := r.Cookie("Authorization")
	if err != nil {
		return "", "", false
	}

	raw, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return "", "", false
	}

	rawParts := strings.Split(string(raw), ":")
	if len(rawParts) != 2 {
		return "", "", false
	}

	return rawParts[0], rawParts[1], true
}
