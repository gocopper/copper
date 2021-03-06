// Package copper provides the primitives to create a new app using github.com/uber-go/fx.
package copper

import (
	"github.com/tusharsoni/copper/chttp"
	"github.com/tusharsoni/copper/clogger"
	"github.com/tusharsoni/copper/crandom"
	"go.uber.org/fx"
)

// NewHTTPApp creates a new copper app that starts a http server.
// It accepts additional modules as fx.Option that can be registered in the app.
// Returns *fx.App that be started using the Run() method.
func NewHTTPApp(opts ...fx.Option) *fx.App {
	combined := append([]fx.Option{
		chttp.Fx,

		fx.Invoke(crandom.Seed),
		fx.Invoke(chttp.Register),
	}, opts...)

	return fx.New(combined...)
}

type HTTPAppParams struct {
	Logger            clogger.Logger
	Routes            []chttp.Route
	GlobalMiddlewares []chttp.MiddlewareFunc
	Config            chttp.Config
}

func RunHTTPApp(p HTTPAppParams) {
	var (
		router = chttp.NewRouter(chttp.RouterParams{
			Routes: append(
				p.Routes,
				chttp.NewHealthRoute(chttp.HealthRouteParams{}).Route,
			),
			GlobalMiddlewareFuncs: append(
				p.GlobalMiddlewares,
				chttp.NewRequestLogger(p.Logger).GlobalMiddlewareFunc,
			),
			Logger: p.Logger,
		})
		server = chttp.NewServer(chttp.ServerParams{
			Logger: p.Logger,
			Config: p.Config,
		})
	)

	crandom.Seed()

	err := chttp.Register(server, router)
	if err != nil {
		p.Logger.Error("Failed to register http router to the server", err)
		return
	}

	err = server.ListenAndServe()
	if err != nil {
		p.Logger.Error("Failed to start the server", err)
		return
	}
}
