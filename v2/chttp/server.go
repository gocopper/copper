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

// NewServerParams holds the params needed to create a server.
type NewServerParams struct {
	Handler http.Handler
	Config  cconfig.Config
	Logger  clogger.Logger
}

// NewServer creates a new server.
func NewServer(p NewServerParams) *Server {
	return &Server{
		handler: p.Handler,
		config:  p.Config,
		logger:  p.Logger,
	}
}

// Server represents a configurable HTTP server that supports graceful shutdown.
type Server struct {
	handler http.Handler
	config  cconfig.Config
	logger  clogger.Logger
}

// Start starts the HTTP server on the configured port and blocks until the context expires.
// Once the context expires, the server is gracefully shutdown.
func (s *Server) Start(ctx context.Context) error {
	var (
		server http.Server
		config config
	)

	err := s.config.Load("chttp", &config)
	if err != nil {
		return cerrors.New(err, "failed to load chttp config", nil)
	}

	server.Addr = fmt.Sprintf(":%d", config.Port)
	server.Handler = s.handler

	go func() {
		s.logger.
			WithTags(map[string]interface{}{"port": config.Port}).
			Info("Starting http server..")

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("Failed to start server", err)
		}
	}()

	<-ctx.Done()

	s.logger.Info("Shutting down http server..")

	ctxShutdown, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(config.ShutdownTimeoutSecs)*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctxShutdown); err != nil {
		s.logger.Error("Failed to shutdown http server cleanly", err)
	}

	return nil
}
