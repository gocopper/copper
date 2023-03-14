package chttptest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gocopper/copper/clogger"

	"github.com/gocopper/copper/chttp"
	"github.com/stretchr/testify/assert"
)

// PingRoutes creates a handler using chttp.NewHandler, starts a test
// http server, and calls each provided route. It verifies that each
// route's handler is called successfully.
func PingRoutes(t *testing.T, routes []chttp.Route) {
	t.Helper()

	for i := range routes {
		body := routes[i].Path
		routes[i].Handler = func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte(body))
			assert.NoError(t, err)
		}
	}

	server := httptest.NewServer(chttp.NewHandler(chttp.NewHandlerParams{
		Routers:           []chttp.Router{NewRouter(routes)},
		GlobalMiddlewares: nil,
		Logger:            clogger.NewNoop(),
	}))
	defer server.Close()

	for _, route := range routes {
		resp, err := http.Get(server.URL + route.Path) //nolint:noctx
		assert.NoError(t, err)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, resp.Body.Close())
		assert.NoError(t, err)

		assert.Equal(t, route.Path, string(body))
	}
}
