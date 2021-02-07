package chttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tusharsoni/copper/v2/clogger"
)

// StartServerParams holds the params to call the StartServer method.
type StartServerParams struct {
	Handler http.Handler
	Config  Config
	Logger  clogger.Logger
	Stop    chan bool
}

// StartServer starts an HTTP server with the given handler.
func StartServer(p StartServerParams) {
	var (
		config = p.Config
		server http.Server
	)

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

	<-p.Stop

	p.Logger.Info("Shutting down http server..")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(p.Config.ShutdownTimeoutSeconds)*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		p.Logger.Error("Failed to shutdown http server cleanly", err)
	}
}

// NewOSSignalStopChan creates a channel compatible with StartServerParams.Stop.
// A message is published on this channel when the process receives an interrupt
// or a term signal.
func NewOSSignalStopChan() chan bool {
	var (
		in  = make(chan os.Signal)
		out = make(chan bool)
	)

	go func() {
		<-in
		signal.Stop(in)
		out <- true
	}()

	signal.Notify(in, syscall.SIGINT, syscall.SIGTERM)

	return out
}
