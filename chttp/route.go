package chttp

import (
	"net/http"

	"go.uber.org/fx"
)

// MiddlewareFunc can be used to create a middleware that can be used on a route.
type MiddlewareFunc func(http.Handler) http.Handler

// RouteResult can be provided to the application container to register a route when starting the http server.
type RouteResult struct {
	fx.Out

	Route Route `group:"routes"`
}

// Route represents a single path that the http server is accepting requests on.
// The route can be configured with middleware functions.
// Additionally, it can be limited to accept requests on specific http methods.
type Route struct {
	MiddlewareFuncs []MiddlewareFunc
	Path            string
	Methods         []string
	Handler         http.Handler
}
