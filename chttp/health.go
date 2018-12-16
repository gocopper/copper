package chttp

import (
	"net/http"

	"go.uber.org/fx"
)

// healthRouteParams holds the dependencies needed to create the Health route using newHealthRoute.
type healthRouteParams struct {
	fx.In

	Config Config `optional:"true"`
}

// newHealthRoute provides a route that responds OK to signify the health of the server.
func newHealthRoute(p healthRouteParams) RouteResult {
	config := p.Config
	if !config.isValid() {
		config = GetDefaultConfig()
	}

	route := Route{
		MiddlewareFuncs: []MiddlewareFunc{},
		Path:            config.HealthPath,
		Methods:         []string{http.MethodGet},
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("OK"))
		}),
	}
	return RouteResult{Route: route}
}
