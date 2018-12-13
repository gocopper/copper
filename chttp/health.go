package chttp

import (
	"net/http"
)

func NewHealthRoute(config Config) RouteResult {
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
