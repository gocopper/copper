//go:build wireinject
// +build wireinject

package copper

import (
	"github.com/gocopper/copper/cconfig"
	"github.com/gocopper/copper/clifecycle"
	"github.com/gocopper/copper/clogger"
	"github.com/google/wire"
)

// InitApp creates a new Copper app along with its dependencies.
func InitApp() (*App, error) {
	panic(
		wire.Build(
			NewApp,
			NewFlags,
			clifecycle.New,
			cconfig.NewWithKeyOverrides,
			clogger.NewWithConfig,
			clogger.LoadConfig,

			wire.FieldsOf(new(*Flags), "ConfigPath"),
		),
	)
}

// WireModule can be used as part of google/wire setup to include the app's
// lifecycle, config, and logger.
var WireModule = wire.NewSet(
	wire.FieldsOf(new(*App), "Lifecycle", "Config", "Logger"),
)
