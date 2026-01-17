//go:build wireinject
// +build wireinject

package copper

import (
	"github.com/gocopper/copper/cconfig"
	"github.com/gocopper/copper/clifecycle"
	"github.com/gocopper/copper/clogger"
	"github.com/google/wire"
)

func InitApp() (*App, error) {
	panic(
		wire.Build(
			NewApp,
			NewFlags,
			clifecycle.New,
			cconfig.NewWithKeyOverrides,
			clogger.NewCore,
			clogger.LoadConfig,

			wire.FieldsOf(new(*Flags), "ConfigPath", "ConfigOverrides"),
		),
	)
}

var WireModule = wire.NewSet(
	wire.FieldsOf(new(*App), "Config", "Lifecycle"),
	clogger.LoadConfig,
	clogger.New,
)
