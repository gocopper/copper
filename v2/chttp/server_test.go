package chttp_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tusharsoni/copper/v2/cconfig"
	"github.com/tusharsoni/copper/v2/cconfig/cconfigtest"
	"github.com/tusharsoni/copper/v2/chttp"
	"github.com/tusharsoni/copper/v2/clogger"
)

func TestServer_Start(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())

	config, err := cconfig.NewConfig(cconfigtest.SetupDirWithConfigs(t, `
[chttp]
Port = 8999
`, ""), "test")
	assert.NoError(t, err)

	server := chttp.NewServer(chttp.NewServerParams{
		Handler: http.NotFoundHandler(),
		Config:  config,
		Logger:  clogger.New(),
	})

	go func() {
		err = server.Start(ctx)
		assert.NoError(t, err)
	}()

	time.Sleep(50 * time.Millisecond) // wait for server to start

	resp, err := http.Get("http://127.0.0.1:8999") //nolint:noctx
	assert.NoError(t, err)
	assert.NoError(t, resp.Body.Close())

	assert.NoError(t, resp.Body.Close())
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	cancel()

	time.Sleep(50 * time.Millisecond) // wait for server to stop

	_, err = http.Get("http://127.0.0.1:8999") //nolint:noctx,bodyclose
	assert.EqualError(t, err, "Get \"http://127.0.0.1:8999\": dial tcp 127.0.0.1:8999: connect: connection refused")
}
