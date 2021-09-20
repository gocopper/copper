package chttp_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gocopper/copper/clogger"

	"github.com/gocopper/copper/chttp"
	"github.com/gocopper/copper/chttp/chttptest"
	"github.com/stretchr/testify/assert"
)

func TestPanicLoggerMiddleware_PanicError(t *testing.T) {
	t.Parallel()

	var (
		router = chttptest.NewRouter([]chttp.Route{
			{
				Path:    "/",
				Methods: []string{http.MethodGet},
				Handler: func(w http.ResponseWriter, r *http.Request) {
					panic(errors.New("test-error"))
				},
			},
		})

		logs = make([]clogger.RecordedLog, 0)

		handler = chttp.NewHandler(chttp.NewHandlerParams{
			Routers: []chttp.Router{router},
			Logger:  clogger.NewRecorder(&logs),
		})
	)

	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL) //nolint:noctx
	assert.NoError(t, err)
	assert.NoError(t, resp.Body.Close())

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, 1, len(logs))
	assert.Equal(t, "Recovered from a panic while handling HTTP request", logs[0].Msg)
	assert.Equal(t, clogger.LevelError, logs[0].Level)
	assert.Equal(t, errors.New("test-error"), logs[0].Error)
}

func TestPanicLoggerMiddleware_NoPanic(t *testing.T) {
	t.Parallel()

	var (
		router = chttptest.NewRouter([]chttp.Route{
			{
				Path:    "/",
				Methods: []string{http.MethodGet},
				Handler: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				},
			},
		})

		logs = make([]clogger.RecordedLog, 0)

		handler = chttp.NewHandler(chttp.NewHandlerParams{
			Routers: []chttp.Router{router},
			Logger:  clogger.NewRecorder(&logs),
		})
	)

	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL) //nolint:noctx
	assert.NoError(t, err)
	assert.NoError(t, resp.Body.Close())

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 0, len(logs))
}

func TestPanicLoggerMiddleware_PanicNonError(t *testing.T) {
	t.Parallel()

	var (
		router = chttptest.NewRouter([]chttp.Route{
			{
				Path:    "/",
				Methods: []string{http.MethodGet},
				Handler: func(w http.ResponseWriter, r *http.Request) {
					panic("test-error")
				},
			},
		})

		logs = make([]clogger.RecordedLog, 0)

		handler = chttp.NewHandler(chttp.NewHandlerParams{
			Routers: []chttp.Router{router},
			Logger:  clogger.NewRecorder(&logs),
		})
	)

	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL) //nolint:noctx
	assert.NoError(t, err)
	assert.NoError(t, resp.Body.Close())

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, 1, len(logs))
	assert.Equal(t, "Recovered from a panic while handling HTTP request", logs[0].Msg)
	assert.Equal(t, clogger.LevelError, logs[0].Level)
	assert.Equal(t, "test-error", logs[0].Tags["error"])
	assert.Nil(t, logs[0].Error)
}
