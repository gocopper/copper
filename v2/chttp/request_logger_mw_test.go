package chttp_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/tusharsoni/copper/v2/chttp"
	"github.com/tusharsoni/copper/v2/chttp/chttptest"
	"github.com/tusharsoni/copper/v2/clogger"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRequestLoggerMiddleware(t *testing.T) {
	t.Parallel()

	var (
		logs   = make([]clogger.RecordedLog, 0)
		logger = clogger.NewRecorder(&logs)
		router = chttptest.NewRouter([]chttp.Route{
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
		})
		handler = chttp.NewHandler(chttp.NewHandlerParams{
			Routers:           []chttp.Router{router},
			GlobalMiddlewares: nil,
		})
	)

	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL + "/test") //nolint:noctx
	assert.NoError(t, err)

	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.NoError(t, resp.Body.Close())

	assert.Equal(t, "OK", string(body))
	assert.Equal(t, clogger.LevelInfo, logs[0].Level)
	assert.Equal(t, "GET /test 201", logs[0].Msg)
}
