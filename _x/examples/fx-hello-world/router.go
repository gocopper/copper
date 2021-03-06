package main

import (
	"net/http"

	"github.com/tusharsoni/copper/chttp"
	"github.com/tusharsoni/copper/clogger"
	"go.uber.org/fx"
)

type Router struct {
	rw     chttp.ReaderWriter
	logger clogger.Logger
}

var RouterFx = fx.Provide(
	NewRouter,
	NewHelloWorldRoute,
)

func NewRouter(rw chttp.ReaderWriter, logger clogger.Logger) *Router {
	return &Router{
		rw:     rw,
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
	ro.rw.OK(w, map[string]string{
		"response": "Hello, World!",
	})
}
