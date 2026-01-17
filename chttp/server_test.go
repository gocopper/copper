package chttp_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/gocopper/copper/chttp"
	"github.com/gocopper/copper/clifecycle/clifecycletest"
	"github.com/gocopper/copper/clogger"
	"github.com/stretchr/testify/assert"
)

func TestServer_Run(t *testing.T) {
	t.Parallel()

	lc := clifecycletest.New()

	server := chttp.NewServer(chttp.NewServerParams{
		Handler:   http.NotFoundHandler(),
		Config:    chttp.Config{Port: 8999},
		Logger:    clogger.NewNoop(),
		Lifecycle: lc,
	})

	go func() {
		err := server.Run()
		assert.NoError(t, err)
	}()

	time.Sleep(50 * time.Millisecond) // wait for server to start

	resp, err := http.Get("http://127.0.0.1:8999") //nolint:noctx
	assert.NoError(t, err)
	assert.NoError(t, resp.Body.Close())

	assert.NoError(t, resp.Body.Close())
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	lc.Stop(clogger.NewNoop())

	time.Sleep(50 * time.Millisecond) // wait for server to stop

	_, err = http.Get("http://127.0.0.1:8999") //nolint:noctx,bodyclose
	assert.EqualError(t, err, "Get \"http://127.0.0.1:8999\": dial tcp 127.0.0.1:8999: connect: connection refused")
}
