package chttp

import (
	"net/http"

	"go.uber.org/fx"
)

type HealthRouteParams struct {
	fx.In

	Config Config `optional:"true"`
}

func NewHealthRoute(p HealthRouteParams) RouteResult {
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
