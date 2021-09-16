package chttp_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gocopper/copper/clogger"

	"github.com/gocopper/copper/chttp"
	"github.com/gocopper/copper/chttp/chttptest"
	"github.com/stretchr/testify/assert"
)

func TestRoutePathInCtxMiddleware(t *testing.T) {
	t.Parallel()

	var (
		routeMWRawRoutePath  string
		globalMWRawRoutePath string

		router = chttptest.NewRouter([]chttp.Route{
			{
				Path:    "/foo/{id}",
				Methods: []string{http.MethodGet},
				Handler: func(w http.ResponseWriter, r *http.Request) {
					routeMWRawRoutePath = chttp.RawRoutePath(r)
				},
			},
		})

		handler = chttp.NewHandler(chttp.NewHandlerParams{
			Routers: []chttp.Router{router},
			GlobalMiddlewares: []chttp.Middleware{
				chttp.HandleMiddleware(func(next http.Handler) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						globalMWRawRoutePath = chttp.RawRoutePath(r)

						next.ServeHTTP(w, r)
					})
				}),
			},
			Logger: clogger.NewNoop(),
		})
	)

	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL + "/foo/bar") //nolint:noctx
	assert.NoError(t, err)
	assert.NoError(t, resp.Body.Close())

	assert.Equal(t, "/foo/{id}", globalMWRawRoutePath)
	assert.Equal(t, "/foo/{id}", routeMWRawRoutePath)
}
