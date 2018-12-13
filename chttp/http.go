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
		p.Logger.Warnw("No routes to register", nil)
	}

	for _, route := range p.Routes {
		p.Logger.Infow("Registering route..", map[string]string{
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

func NewServer(lc fx.Lifecycle, config Config) *http.ServeMux {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: mux,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go server.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})

	return mux
}
