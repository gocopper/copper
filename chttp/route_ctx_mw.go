package chttp

import (
	"context"
	"net/http"
)

type ctxRoutePath string

const ctxRoutePathKey = ctxRoutePath("chttp/route-path")

func setRoutePathInCtxMiddleware(path string) Middleware {
	var mw = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), ctxRoutePathKey, path)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	return HandleMiddleware(mw)
}

// RawRoutePath returns the route path that matched for the given http.Request. This path includes the raw URL
// variables. For example, a route path "/foo/{id}" will be returned as-is (i.e. {id} will NOT be replaced with the
// actual url path)
func RawRoutePath(r *http.Request) string {
	return r.Context().Value(ctxRoutePathKey).(string)
}
