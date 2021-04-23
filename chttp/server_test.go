package chttp_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tusharsoni/copper"
	"github.com/tusharsoni/copper/cconfig"
	"github.com/tusharsoni/copper/cconfig/cconfigtest"
	"github.com/tusharsoni/copper/chttp"
	"github.com/tusharsoni/copper/clogger"
)

func TestServer_Run(t *testing.T) {
	t.Parallel()

	logger := clogger.New()
	lc := copper.NewLifecycle(logger)

	config, err := cconfig.New(cconfigtest.SetupDirWithConfigs(t, `
[chttp]
Port = 8999
`, ""), ".", "test")
	assert.NoError(t, err)

	server := chttp.NewServer(chttp.NewServerParams{
		Handler:   http.NotFoundHandler(),
		Config:    config,
		Logger:    logger,
		Lifecycle: lc,
	})

	go func() {
		err = server.Run()
		assert.NoError(t, err)
	}()

	time.Sleep(50 * time.Millisecond) // wait for server to start

	resp, err := http.Get("http://127.0.0.1:8999") //nolint:noctx
	assert.NoError(t, err)
	assert.NoError(t, resp.Body.Close())

	assert.NoError(t, resp.Body.Close())
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	lc.Stop()

	time.Sleep(50 * time.Millisecond) // wait for server to stop

	_, err = http.Get("http://127.0.0.1:8999") //nolint:noctx,bodyclose
	assert.EqualError(t, err, "Get \"http://127.0.0.1:8999\": dial tcp 127.0.0.1:8999: connect: connection refused")
}
