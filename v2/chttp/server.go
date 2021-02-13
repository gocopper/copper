package chttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/tusharsoni/copper/v2/cconfig"
	"github.com/tusharsoni/copper/v2/cerrors"
	"github.com/tusharsoni/copper/v2/clogger"
)

// StartServerParams holds the params to call the StartServer method.
type StartServerParams struct {
	Handler http.Handler
	Config  cconfig.Config
	Logger  clogger.Logger
}

// StartServer starts an HTTP server with the given handler.
func StartServer(ctx context.Context, p StartServerParams) error {
	var (
		server http.Server
		config config
	)

	err := p.Config.Load("chttp", &config)
	if err != nil {
		return cerrors.New(err, "failed to load chttp config", nil)
	}

	server.Addr = fmt.Sprintf(":%d", config.Port)
	server.Handler = p.Handler

	go func() {
		p.Logger.
			WithTags(map[string]interface{}{"port": config.Port}).
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
		time.Duration(config.ShutdownTimeoutSecs)*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctxShutdown); err != nil {
		p.Logger.Error("Failed to shutdown http server cleanly", err)
	}

	return nil
}
