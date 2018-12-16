// Package copper provides the primitives to create a new app using github.com/uber-go/fx.
package copper

import (
	"github.com/tusharsoni/copper/chttp"

	"go.uber.org/fx"
)

// NewApp creates a new copper app that starts a http server.
// It accepts additional modules as fx.Option that can be registered in the app.
// Returns *fx.App that be started using the Run() method.
func NewApp(opts ...fx.Option) *fx.App {
	combined := append([]fx.Option{
		chttp.Fx,

		fx.Invoke(chttp.Register),
	}, opts...)

	return fx.New(combined...)
}
