package chttp

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/tusharsoni/copper/clogger"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

// Register can be used with fx.Invoke to register the root handler with the server.
func Register(server *http.ServeMux, handler http.Handler) {
	server.Handle("/", handler)
}

// routerParams holds the dependencies needed to create a router using newRouter.
type routerParams struct {
	fx.In

	Routes []Route `group:"routes"`
	Logger clogger.Logger
}

// newRouter creates a http.Handler by registering all routes that have been provided in the application container.
func newRouter(p routerParams) http.Handler {
	r := mux.NewRouter()

	if len(p.Routes) == 0 {
		p.Logger.Warn("No routes to register", nil)
	}

	for _, route := range p.Routes {
		p.Logger.Info("Registering route..", map[string]string{
			"path":    route.Path,
			"methods": strings.Join(route.Methods, ", "),
		})

		handlerFunc := route.Handler

		for _, f := range route.MiddlewareFuncs {
			handlerFunc = f(route.Handler)
		}

		muxRoute := r.Handle(route.Path, handlerFunc)

		if len(route.Methods) > 0 {
			muxRoute.Methods(route.Methods...)
		}
	}

	return r
}

// serverParams holds the dependencies needed to create a http server using newServer.
type serverParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    clogger.Logger

	Config Config `optional:"true"`
}

// newServer creates a http request mux. This server starts when the application starts and stops gracefully when
// the application stops.
func newServer(p serverParams) *http.ServeMux {
	config := p.Config
	if !config.isValid() {
		config = GetDefaultConfig()
	}

	serveMux := http.NewServeMux()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: serveMux,
	}

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

	return serveMux
}
