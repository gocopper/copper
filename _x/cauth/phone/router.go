package phone

import (
	"net/http"

	"github.com/tusharsoni/copper/chttp"
	"github.com/tusharsoni/copper/clogger"
	"go.uber.org/fx"
)

type Router struct {
	rw     chttp.ReaderWriter
	logger clogger.Logger

	auth Svc
}

type RouterParams struct {
	fx.In

	RW     chttp.ReaderWriter
	Logger clogger.Logger

	Auth Svc
}

func NewRouter(p RouterParams) *Router {
	return &Router{
		rw:     p.RW,
		logger: p.Logger,

		auth: p.Auth,
	}
}

func NewSignup(ro *Router) chttp.RouteResult {
	return chttp.RouteResult{Route: chttp.Route{
		Path:            "/api/auth/phone/signup",
		MiddlewareFuncs: []chttp.MiddlewareFunc{},
		Methods:         []string{http.MethodPost},
		Handler:         http.HandlerFunc(ro.HandleSignup),
	}}
}

func (ro *Router) HandleSignup(w http.ResponseWriter, r *http.Request) {
	var body struct {
		PhoneNumber string `json:"phone_number" valid:"auth.PhoneNumber,required"`
	}

	if !ro.rw.Read(w, r, &body) {
		return
	}

	err := ro.auth.Signup(r.Context(), body.PhoneNumber)
	if err != nil {
		ro.logger.Error("Failed to sign up with phone number", err)
		ro.rw.InternalErr(w)
		return
	}

	ro.rw.OK(w, nil)
}

func NewLogin(ro *Router) chttp.RouteResult {
	return chttp.RouteResult{Route: chttp.Route{
		Path:            "/api/auth/phone/login",
		MiddlewareFuncs: []chttp.MiddlewareFunc{},
		Methods:         []string{http.MethodPost},
		Handler:         http.HandlerFunc(ro.HandleLogin),
	}}
}

func (ro *Router) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var body struct {
		PhoneNumber      string `json:"phone_number" valid:"auth.PhoneNumber,required"`
		VerificationCode uint   `json:"verification_code" valid:"required"`
	}

	if !ro.rw.Read(w, r, &body) {
		return
	}

	c, sessionToken, err := ro.auth.Login(r.Context(), body.PhoneNumber, body.VerificationCode)
	if err != nil {
		ro.logger.WithTags(map[string]interface{}{
			"phoneNumber": body.PhoneNumber,
		}).Error("Failed to login with phone number", err)
		ro.rw.InternalErr(w)
		return
	}

	ro.rw.OK(w, &struct {
		UserUUID     string `json:"user_uuid"`
		SessionToken string `json:"session_token"`
	}{
		UserUUID:     c.UserUUID,
		SessionToken: sessionToken,
	})
}
