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

func Register(server *http.ServeMux, handler http.Handler) {
	server.Handle("/", handler)
}

type RouterParams struct {
	fx.In

	Routes []Route `group:"routes"`
	Logger clogger.Logger
}

func NewRouter(p RouterParams) http.Handler {
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

func NewServer(lc fx.Lifecycle, logger clogger.Logger, config Config) *http.ServeMux {
	serveMux := http.NewServeMux()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: serveMux,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				err := server.ListenAndServe()
				if err != nil && err != http.ErrServerClosed {
					logger.Error("Failed to start http server", err)
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
