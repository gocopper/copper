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
