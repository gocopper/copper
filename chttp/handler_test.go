package chttp_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gocopper/copper/chttp"
	"github.com/gocopper/copper/chttp/chttptest"
	"github.com/stretchr/testify/assert"
)

func TestNewHandler(t *testing.T) {
	t.Parallel()

	router := chttptest.NewRouter([]chttp.Route{
		{
			Path:    "/",
			Methods: []string{http.MethodGet},
			Handler: func(w http.ResponseWriter, r *http.Request) {
				_, err := w.Write([]byte("success"))
				assert.NoError(t, err)
			},
		},
	})

	server := httptest.NewServer(chttp.NewHandler(chttp.NewHandlerParams{
		Routers:           []chttp.Router{router},
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

	router := chttptest.NewRouter([]chttp.Route{
		{
			Path:    "/",
			Handler: func(w http.ResponseWriter, r *http.Request) {},
		},
	})

	server := httptest.NewServer(chttp.NewHandler(chttp.NewHandlerParams{
		Routers: []chttp.Router{router},
		GlobalMiddlewares: []chttp.Middleware{
			chttptest.NewMiddleware(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					didCallGlobalMiddleware = true
					next.ServeHTTP(w, r)
				})
			}),
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

	middlewareCalls := make([]string, 0)

	router := chttptest.NewRouter([]chttp.Route{
		{
			Path: "/",
			Middlewares: []chttp.Middleware{
				chttptest.NewMiddleware(func(next http.Handler) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						middlewareCalls = append(middlewareCalls, "1")
						next.ServeHTTP(w, r)
					})
				}),

				chttptest.NewMiddleware(func(next http.Handler) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						middlewareCalls = append(middlewareCalls, "2")
						next.ServeHTTP(w, r)
					})
				}),
			},
			Handler: func(w http.ResponseWriter, r *http.Request) {},
		},
	})

	server := httptest.NewServer(chttp.NewHandler(chttp.NewHandlerParams{
		Routers:           []chttp.Router{router},
		GlobalMiddlewares: nil,
	}))
	defer server.Close()

	resp, err := http.Get(server.URL) //nolint:noctx
	assert.NoError(t, err)
	assert.NoError(t, resp.Body.Close())

	assert.Equal(t, []string{"1", "2"}, middlewareCalls)
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
