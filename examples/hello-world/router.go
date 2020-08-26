package main

import (
	"net/http"

	"github.com/tusharsoni/copper/chttp"
)

type Router struct {
	rw chttp.ReaderWriter
}

func NewRouter(rw chttp.ReaderWriter) *Router {
	return &Router{
		rw: rw,
	}
}

func (ro *Router) Routes() []chttp.Route {
	return []chttp.Route{
		{
			Path:    "/",
			Methods: []string{http.MethodGet},
			Handler: http.HandlerFunc(ro.HandleHelloWorld),
		},
	}
}

func (ro *Router) HandleHelloWorld(w http.ResponseWriter, r *http.Request) {

	ro.rw.OK(w, map[string]string{
		"response": "Hello, World",
	})
}
