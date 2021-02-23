package chttptest

import "github.com/tusharsoni/copper/v2/chttp"

// ReverseRoutes reverses the provided slice of chttp.Route.
func ReverseRoutes(routes []chttp.Route) []chttp.Route {
	for i := 0; i < len(routes)/2; i++ {
		j := len(routes) - i - 1
		routes[i], routes[j] = routes[j], routes[i]
	}

	return routes
}

// NewRouter returns a router that returns the given routes.
func NewRouter(routes []chttp.Route) chttp.Router {
	return &router{routes: routes}
}

type router struct {
	routes []chttp.Route
}

func (ro *router) Routes() []chttp.Route {
	return ro.routes
}
