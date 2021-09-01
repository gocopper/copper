package chttp_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gocopper/copper/chttp"
	"github.com/gocopper/copper/chttp/chttptest"
	"github.com/stretchr/testify/assert"
)

func TestRoutePathInCtxMiddleware(t *testing.T) {
	t.Parallel()

	var (
		rawRoutePath string

		router = chttptest.NewRouter([]chttp.Route{
			{
				Path:    "/foo/{id}",
				Methods: []string{http.MethodGet},
				Handler: func(w http.ResponseWriter, r *http.Request) {
					rawRoutePath = chttp.RawRoutePath(r)
				},
			},
		})

		handler = chttp.NewHandler(chttp.NewHandlerParams{
			Routers: []chttp.Router{router},
		})
	)

	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL + "/foo/bar") //nolint:noctx
	assert.NoError(t, err)
	assert.NoError(t, resp.Body.Close())

	assert.Equal(t, "/foo/{id}", rawRoutePath)
}
