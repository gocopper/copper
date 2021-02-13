package chttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/tusharsoni/copper/v2/cconfig"
	"github.com/tusharsoni/copper/v2/clogger"
)

const (
	defaultPort                = 7501
	defaultShutdownTimeoutSecs = 15
)

// StartServerParams holds the params to call the StartServer method.
type StartServerParams struct {
	Handler http.Handler
	Config  cconfig.Config
	Logger  clogger.Logger
}

// StartServer starts an HTTP server with the given handler.
func StartServer(ctx context.Context, p StartServerParams) {
	var server http.Server

	port, ok := p.Config.Value("chttp.port").(int64)
	if !ok {
		port = defaultPort
	}

	shutdownTimeoutSecs, ok := p.Config.Value("chttp.shutdown_timeout_secs").(int64)
	if !ok {
		shutdownTimeoutSecs = defaultShutdownTimeoutSecs
	}

	server.Addr = fmt.Sprintf(":%d", port)
	server.Handler = p.Handler

	go func() {
		p.Logger.
			WithTags(map[string]interface{}{"port": port}).
			Info("Starting http server..")

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			p.Logger.Error("Failed to start server", err)
		}
	}()

	<-ctx.Done()

	p.Logger.Info("Shutting down http server..")

	ctxShutdown, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(shutdownTimeoutSecs)*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctxShutdown); err != nil {
		p.Logger.Error("Failed to shutdown http server cleanly", err)
	}
}
