package chttp

import (
	"net/http"

	"go.uber.org/fx"
)

type MiddlewareFunc func(http.Handler) http.Handler

type RouteResult struct {
	fx.Out

	Route Route `group:"routes"`
}

type Route struct {
	MiddlewareFuncs []MiddlewareFunc
	Path            string
	Methods         []string
	Handler         http.Handler
}
