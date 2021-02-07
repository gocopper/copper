package chttp_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tusharsoni/copper/v2/chttp"
	"github.com/tusharsoni/copper/v2/chttp/chttptest"
)

func TestNewHandler(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(chttp.NewHandler(chttp.NewHandlerParams{
		Routes: []chttp.Route{
			{
				Path:    "/",
				Methods: []string{http.MethodGet},
				Handler: func(w http.ResponseWriter, r *http.Request) {
					_, err := w.Write([]byte("success"))
					assert.NoError(t, err)
				},
			},
		},
		GlobalMiddlewares: nil,
	}))
	defer server.Close()

	resp, err := http.Get(server.URL) //nolint:noctx
	assert.NoError(t, err)

	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, resp.Body.Close())
	assert.NoError(t, err)

	assert.Equal(t, "success", string(body))
}

func TestNewHandler_GlobalMiddleware(t *testing.T) {
	t.Parallel()

	didCallGlobalMiddleware := false

	server := httptest.NewServer(chttp.NewHandler(chttp.NewHandlerParams{
		Routes: []chttp.Route{
			{
				Path:    "/",
				Handler: func(w http.ResponseWriter, r *http.Request) {},
			},
		},
		GlobalMiddlewares: []chttp.Middleware{
			func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					didCallGlobalMiddleware = true
					next.ServeHTTP(w, r)
				})
			},
		},
	}))
	defer server.Close()

	resp, err := http.Get(server.URL) //nolint:noctx
	assert.NoError(t, err)
	assert.NoError(t, resp.Body.Close())

	assert.True(t, didCallGlobalMiddleware)
}

func TestNewHandler_RouteMiddleware(t *testing.T) {
	t.Parallel()

	didCallRouteMiddleware := false

	server := httptest.NewServer(chttp.NewHandler(chttp.NewHandlerParams{
		Routes: []chttp.Route{
			{
				Path: "/",
				Middlewares: []chttp.Middleware{
					func(next http.Handler) http.Handler {
						return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							didCallRouteMiddleware = true
							next.ServeHTTP(w, r)
						})
					},
				},
				Handler: func(w http.ResponseWriter, r *http.Request) {},
			},
		},
		GlobalMiddlewares: nil,
	}))
	defer server.Close()

	resp, err := http.Get(server.URL) //nolint:noctx
	assert.NoError(t, err)
	assert.NoError(t, resp.Body.Close())

	assert.True(t, didCallRouteMiddleware)
}

func TestNewHandler_RoutePriority_WithPlaceholder(t *testing.T) {
	t.Parallel()

	routes := []chttp.Route{
		{Path: "/foo"},
		{Path: "/{id}"},
	}

	chttptest.PingRoutes(t, routes)
	chttptest.PingRoutes(t, chttptest.ReverseRoutes(routes))
}

func TestNewHandler_RoutePriority_WithIndex(t *testing.T) {
	t.Parallel()

	routes := []chttp.Route{
		{Path: "/foo"},
		{Path: "/"},
	}

	chttptest.PingRoutes(t, routes)
	chttptest.PingRoutes(t, chttptest.ReverseRoutes(routes))
}

func TestNewHandler_RoutePriority_Equal(t *testing.T) {
	t.Parallel()

	routes := []chttp.Route{
		{Path: "/foo"},
		{Path: "/bar"},
	}

	chttptest.PingRoutes(t, routes)
	chttptest.PingRoutes(t, chttptest.ReverseRoutes(routes))
}
