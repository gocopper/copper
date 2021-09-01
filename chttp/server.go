package chttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gocopper/copper"
	"github.com/gocopper/copper/clogger"
)

// NewServerParams holds the params needed to create a server.
type NewServerParams struct {
	Handler   http.Handler
	Lifecycle *copper.Lifecycle
	Config    Config
	Logger    clogger.Logger
}

// NewServer creates a new server.
func NewServer(p NewServerParams) *Server {
	return &Server{
		handler:  p.Handler,
		config:   p.Config,
		logger:   p.Logger,
		lc:       p.Lifecycle,
		internal: http.Server{},
	}
}

// Server represents a configurable HTTP server that supports graceful shutdown.
type Server struct {
	handler http.Handler
	config  Config
	logger  clogger.Logger
	lc      *copper.Lifecycle

	internal http.Server
}

// Run configures an HTTP server using the provided app config and starts it.
func (s *Server) Run() error {
	s.internal.Addr = fmt.Sprintf(":%d", s.config.Port)
	s.internal.Handler = s.handler

	s.lc.OnStop(func(ctx context.Context) error {
		s.logger.Info("Shutting down http server..")

		return s.internal.Shutdown(ctx)
	})

	go func() {
		s.logger.
			WithTags(map[string]interface{}{"port": s.config.Port}).
			Info("Starting http server..")

		err := s.internal.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("Server did not close cleanly", err)
		}
	}()

	return nil
}
