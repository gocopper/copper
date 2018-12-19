package cauth

import (
	"net/http"

	"github.com/tusharsoni/copper/clogger"

	"github.com/tusharsoni/copper/chttp"
)

type router struct {
	req    chttp.BodyReader
	resp   chttp.Responder
	users  UsersSvc
	logger clogger.Logger
}

func newRouter(req chttp.BodyReader, resp chttp.Responder, users UsersSvc, logger clogger.Logger) *router {
	return &router{
		req:    req,
		resp:   resp,
		users:  users,
		logger: logger,
	}
}

func newLoginRoute(ro *router) chttp.RouteResult {
	route := chttp.Route{
		Path:    "/login",
		Methods: []string{http.MethodPost},
		Handler: http.HandlerFunc(ro.login),
	}
	return chttp.RouteResult{Route: route}
}

func (ro *router) login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email" valid:"email"`
		Password string `json:"password" valid:"runelength(4|32)"`
	}
	var resp struct {
		User         *user  `json:"user"`
		SessionToken string `json:"session_token"`
	}

	if !ro.req.Read(w, r, &body) {
		return
	}

	u, sessionToken, err := ro.users.Login(r.Context(), body.Email, body.Password)
	if err != nil && err != ErrInvalidCredentials {
		ro.logger.Error("Failed to login user with email and password", err)
		ro.resp.InternalErr(w)
		return
	} else if err == ErrInvalidCredentials {
		ro.resp.Unauthorized(w)
		return
	}

	resp.User = u
	resp.SessionToken = sessionToken

	ro.resp.OK(w, resp)
}

func newSignupRoute(ro *router) chttp.RouteResult {
	route := chttp.Route{
		Path:    "/signup",
		Methods: []string{http.MethodPost},
		Handler: http.HandlerFunc(ro.signup),
	}
	return chttp.RouteResult{Route: route}
}

func (ro *router) signup(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email" valid:"email"`
		Password string `json:"password" valid:"runelength(4|32)"`
	}

	if !ro.req.Read(w, r, &body) {
		return
	}

	user, err := ro.users.Signup(r.Context(), body.Email, body.Password)
	if err != nil && err != ErrUserAlreadyExists {
		ro.logger.Error("Failed to signup user with email and password", err)
		ro.resp.InternalErr(w)
		return
	} else if err == ErrUserAlreadyExists {
		ro.resp.BadRequest(w, ErrUserAlreadyExists)
		return
	}

	ro.resp.Created(w, user)
}
