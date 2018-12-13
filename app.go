package copper

import (
	"github.com/tusharsoni/copper/chttp"

	"go.uber.org/fx"
)

func NewApp(opts ...fx.Option) *fx.App {
	combined := append([]fx.Option{
		chttp.Fx,

		fx.Invoke(chttp.Register),
	}, opts...)

	return fx.New(combined...)
}
