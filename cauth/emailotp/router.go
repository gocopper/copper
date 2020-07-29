package emailotp

import (
	"net/http"

	"github.com/tusharsoni/copper/chttp"
	"github.com/tusharsoni/copper/clogger"
	"go.uber.org/fx"
)

type Router struct {
	resp   chttp.Responder
	req    chttp.BodyReader
	logger clogger.Logger

	auth Svc
}

type RouterParams struct {
	fx.In

	Resp   chttp.Responder
	Req    chttp.BodyReader
	Logger clogger.Logger

	Auth Svc
}

func NewRouter(p RouterParams) *Router {
	return &Router{
		resp:   p.Resp,
		req:    p.Req,
		logger: p.Logger,

		auth: p.Auth,
	}
}

func NewSignup(ro *Router) chttp.RouteResult {
	return chttp.RouteResult{Route: chttp.Route{
		Path:    "/api/auth/email-otp/signup",
		Methods: []string{http.MethodPost},
		Handler: http.HandlerFunc(ro.HandleSignup),
	}}
}

func (ro *Router) HandleSignup(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email string `json:"email" valid:"email,required"`
	}

	if !ro.req.Read(w, r, &body) {
		return
	}

	c, token, err := ro.auth.Signup(r.Context(), body.Email)
	if err != nil {
		ro.logger.Error("Failed to sign up with email", err)
		ro.resp.InternalErr(w)
		return
	}

	ro.resp.OK(w, map[string]interface{}{
		"user_uuid":     c.UserUUID,
		"session_token": token,
	})
}

func NewLogin(ro *Router) chttp.RouteResult {
	return chttp.RouteResult{Route: chttp.Route{
		Path:    "/api/auth/email-otp/login",
		Methods: []string{http.MethodPost},
		Handler: http.HandlerFunc(ro.HandleLogin),
	}}
}

func (ro *Router) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email            string `json:"email" valid:"email,required"`
		VerificationCode uint   `json:"verification_code" valid:"required"`
	}

	if !ro.req.Read(w, r, &body) {
		return
	}

	c, sessionToken, err := ro.auth.Login(r.Context(), body.Email, body.VerificationCode)
	if err != nil {
		ro.logger.WithTags(map[string]interface{}{
			"email": body.Email,
		}).Error("Failed to login with email", err)
		ro.resp.InternalErr(w)
		return
	}

	ro.resp.OK(w, map[string]string{
		"user_uuid":     c.UserUUID,
		"session_token": sessionToken,
	})
}
