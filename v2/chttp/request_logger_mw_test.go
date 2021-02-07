package chttp_test

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tusharsoni/copper/v2/chttp"
	"github.com/tusharsoni/copper/v2/clogger/console"
)

func TestNewRequestLoggerMiddleware(t *testing.T) {
	t.Parallel()

	var (
		buf     bytes.Buffer
		logger  = console.New()
		handler = chttp.NewHandler(chttp.NewHandlerParams{
			Routes: []chttp.Route{
				{
					Middlewares: []chttp.Middleware{
						chttp.NewRequestLoggerMiddleware(logger),
					},
					Path:    "/test",
					Methods: []string{http.MethodGet},
					Handler: func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(201)

						_, err := w.Write([]byte("OK"))
						assert.NoError(t, err)
					},
				},
			},
			GlobalMiddlewares: nil,
		})
	)

	log.SetOutput(&buf)

	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL + "/test") //nolint:noctx
	assert.NoError(t, err)

	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.NoError(t, resp.Body.Close())

	assert.Equal(t, "OK", string(body))
	assert.Contains(t, buf.String(), "[INFO] GET /test 201")
}
