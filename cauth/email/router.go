package email

import (
	"net/http"

	"github.com/tusharsoni/copper/cauth"

	"github.com/tusharsoni/copper/chttp"
	"github.com/tusharsoni/copper/clogger"
	"go.uber.org/fx"
)

type Router struct {
	resp   chttp.Responder
	req    chttp.BodyReader
	logger clogger.Logger

	auth   Svc
	config Config
}

type RouterParams struct {
	fx.In

	Resp   chttp.Responder
	Req    chttp.BodyReader
	Logger clogger.Logger

	Auth   Svc
	Config Config
}

func NewRouter(p RouterParams) *Router {
	return &Router{
		resp:   p.Resp,
		req:    p.Req,
		logger: p.Logger,

		auth:   p.Auth,
		config: p.Config,
	}
}

func NewSignupRoute(ro *Router) chttp.RouteResult {
	route := chttp.Route{
		Path:    "/api/auth/email/signup",
		Methods: []string{http.MethodPost},
		Handler: http.HandlerFunc(ro.Signup),
	}
	return chttp.RouteResult{Route: route}
}

func (ro *Router) Signup(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email" valid:"email"`
		Password string `json:"password" valid:"runelength(4|32)"`
	}

	if !ro.req.Read(w, r, &body) {
		return
	}

	c, sessionToken, err := ro.auth.Signup(r.Context(), body.Email, body.Password)
	if err != nil && err != cauth.ErrUserAlreadyExists {
		ro.logger.Error("Failed to signup user with email and password", err)
		ro.resp.InternalErr(w)
		return
	} else if err == cauth.ErrUserAlreadyExists {
		ro.resp.BadRequest(w, cauth.ErrUserAlreadyExists)
		return
	}

	ro.resp.Created(w, map[string]string{
		"user_uuid":     c.UserUUID,
		"session_token": sessionToken,
	})
}

func NewLoginRoute(ro *Router) chttp.RouteResult {
	route := chttp.Route{
		Path:    "/api/auth/email/login",
		Methods: []string{http.MethodPost},
		Handler: http.HandlerFunc(ro.Login),
	}
	return chttp.RouteResult{Route: route}
}

func (ro *Router) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email" valid:"email"`
		Password string `json:"password" valid:"runelength(4|32)"`
	}

	if !ro.req.Read(w, r, &body) {
		return
	}

	c, sessionToken, err := ro.auth.Login(r.Context(), body.Email, body.Password)
	if err != nil && err != cauth.ErrInvalidCredentials {
		ro.logger.Error("Failed to login user with email and password", err)
		ro.resp.InternalErr(w)
		return
	} else if err == cauth.ErrInvalidCredentials {
		ro.resp.Unauthorized(w)
		return
	}

	ro.resp.OK(w, map[string]string{
		"user_uuid":     c.UserUUID,
		"session_token": sessionToken,
	})
}

func NewVerifyUserRoute(ro *Router, mw cauth.Middleware) chttp.RouteResult {
	route := chttp.Route{
		MiddlewareFuncs: []chttp.MiddlewareFunc{mw.VerifySessionToken},
		Path:            "/api/auth/email/verify",
		Methods:         []string{http.MethodPost},
		Handler:         http.HandlerFunc(ro.VerifyUser),
	}
	return chttp.RouteResult{Route: route}
}

func (ro *Router) VerifyUser(w http.ResponseWriter, r *http.Request) {
	var body struct {
		VerificationCode string `json:"verification_code" valid:"printableascii"`
	}

	if !ro.req.Read(w, r, &body) {
		return
	}

	userUUID := cauth.GetCurrentUserUUID(r.Context())

	err := ro.auth.VerifyUser(r.Context(), userUUID, body.VerificationCode)
	if err != nil && err != cauth.ErrInvalidCredentials {
		ro.logger.Error("Failed to verify user", err)
		ro.resp.InternalErr(w)
		return
	} else if err == cauth.ErrInvalidCredentials {
		ro.resp.BadRequest(w, err)
		return
	}

	ro.resp.OK(w, nil)
}

func NewResendVerificationCodeRoute(ro *Router, mw cauth.Middleware) chttp.RouteResult {
	route := chttp.Route{
		MiddlewareFuncs: []chttp.MiddlewareFunc{mw.VerifySessionToken},
		Path:            "/api/auth/email/resend-verification-code",
		Methods:         []string{http.MethodPost},
		Handler:         http.HandlerFunc(ro.ResendVerificationCode),
	}
	return chttp.RouteResult{Route: route}
}

func (ro *Router) ResendVerificationCode(w http.ResponseWriter, r *http.Request) {
	userUUID := cauth.GetCurrentUserUUID(r.Context())

	err := ro.auth.ResendVerificationCode(r.Context(), userUUID)
	if err != nil {
		ro.logger.Error("Failed to resend verification code", err)
		ro.resp.InternalErr(w)
		return
	}

	ro.resp.OK(w, nil)
}

func NewChangePasswordRoute(ro *Router) chttp.RouteResult {
	route := chttp.Route{
		Path:    "/api/auth/email/change-password",
		Methods: []string{http.MethodPost},
		Handler: http.HandlerFunc(ro.ChangePassword),
	}
	return chttp.RouteResult{Route: route}
}

func (ro *Router) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email       string `json:"email" valid:"email"`
		OldPassword string `json:"old_password" valid:"printableascii"`
		NewPassword string `json:"new_password" valid:"printableascii"`
	}

	if !ro.req.Read(w, r, &body) {
		return
	}

	err := ro.auth.ChangePassword(r.Context(), body.Email, body.OldPassword, body.NewPassword)
	if err != nil {
		ro.logger.Error("Failed to change password", err)
		ro.resp.InternalErr(w)
		return
	}

	ro.resp.OK(w, nil)
}

func NewResetPasswordRoute(ro *Router) chttp.RouteResult {
	route := chttp.Route{
		Path:    "/api/auth/email/reset-password",
		Methods: []string{http.MethodPost},
		Handler: http.HandlerFunc(ro.ResetPassword),
	}
	return chttp.RouteResult{Route: route}
}

func (ro *Router) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email string `json:"email" valid:"email"`
	}

	if !ro.req.Read(w, r, &body) {
		return
	}

	err := ro.auth.ResetPassword(r.Context(), body.Email)
	if err != nil {
		ro.logger.Error("Failed to reset password", err)
		ro.resp.InternalErr(w)
		return
	}

	ro.resp.OK(w, nil)
}
