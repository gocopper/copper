package chttp_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tusharsoni/copper/v2/cconfig"
	"github.com/tusharsoni/copper/v2/chttp"
	"github.com/tusharsoni/copper/v2/clogger"
)

func TestStartServer(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		chttp.StartServer(ctx, chttp.StartServerParams{
			Handler: http.NotFoundHandler(),
			Config: cconfig.NewStaticConfig(map[string]interface{}{
				"chttp.port": int64(8999),
			}),
			Logger: clogger.NewNoop(),
		})
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
