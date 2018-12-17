// Package copper provides the primitives to create a new app using github.com/uber-go/fx.
package copper

import (
	"github.com/tusharsoni/copper/chttp"
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
