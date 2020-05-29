package main

import (
	"net/http"

	"github.com/tusharsoni/copper/chttp"
	"github.com/tusharsoni/copper/clogger"
	"go.uber.org/fx"
)

type Router struct {
	req    chttp.BodyReader
	resp   chttp.Responder
	logger clogger.Logger
}

var RouterFx = fx.Provide(
	NewRouter,
	NewHelloWorldRoute,
)

func NewRouter(req chttp.BodyReader, resp chttp.Responder, logger clogger.Logger) *Router {
	return &Router{
		req:    req,
		resp:   resp,
		logger: logger,
	}
}

func NewHelloWorldRoute(ro *Router) chttp.RouteResult {
	return chttp.RouteResult{Route: chttp.Route{
		Path:    "/",
		Methods: []string{http.MethodGet},
		Handler: http.HandlerFunc(ro.HandleHelloWorld),
	}}
}

func (ro *Router) HandleHelloWorld(w http.ResponseWriter, r *http.Request) {
	ro.resp.OK(w, map[string]string{
		"response": "Hello, World!",
	})
}
