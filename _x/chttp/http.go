package chttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/tusharsoni/copper/clogger"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

// Register can be used with fx.Invoke to register the root handler with the server.
func Register(server *http.Server, handler http.Handler) error {
	muxHandler, ok := server.Handler.(*http.ServeMux)
	if !ok {
		return errors.New("http server does not have a *http.ServeMux handler")
	}

	muxHandler.Handle("/", handler)

	return nil
}

// RouterParams holds the dependencies needed to create a router using NewRouter.
type RouterParams struct {
	fx.In

	Routes                []Route          `group:"routes"`
	GlobalMiddlewareFuncs []MiddlewareFunc `group:"global_middleware_funcs"`
	Logger                clogger.Logger
}

// NewRouter creates a http.Handler by registering all routes that have been provided in the application container.
func NewRouter(p RouterParams) http.Handler {
	r := mux.NewRouter()

	if len(p.Routes) == 0 {
		p.Logger.Warn("No routes to register", nil)
	}

	for _, f := range p.GlobalMiddlewareFuncs {
		r.Use(mux.MiddlewareFunc(f))
	}

	sortRoutes(p.Routes)

	for _, route := range p.Routes {
		p.Logger.WithTags(map[string]interface{}{
			"path":    route.Path,
			"methods": strings.Join(route.Methods, ", "),
		}).Info("Registering route..")

		handlerFunc := route.Handler

		for _, f := range route.MiddlewareFuncs {
			handlerFunc = f(handlerFunc)
		}

		muxRoute := r.Handle(route.Path, handlerFunc)

		if len(route.Methods) > 0 {
			muxRoute.Methods(route.Methods...)
		}
	}

	return r
}

// ServerParams holds the dependencies needed to create a http server using NewServer.
type ServerParams struct {
	fx.In

	Logger clogger.Logger

	Config    Config       `optional:"true"`
	Lifecycle fx.Lifecycle `optional:"true"`
}

// NewServer creates a http server. This server starts when the application starts and stops gracefully when
// the application stops.
func NewServer(p ServerParams) *http.Server {
	config := p.Config
	if !config.isValid() {
		config = GetDefaultConfig()
	}

	p.Logger.WithTags(map[string]interface{}{
		"port": config.Port,
	}).Info("Setting up HTTP server..")

	serveMux := http.NewServeMux()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: serveMux,
	}

	if p.Lifecycle != nil {
		p.Lifecycle.Append(fx.Hook{
			OnStart: func(context.Context) error {
				go func() {
					err := server.ListenAndServe()
					if err != nil && err != http.ErrServerClosed {
						p.Logger.Error("Failed to start http server", err)
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return server.Shutdown(ctx)
			},
		})
	}

	return server
}
