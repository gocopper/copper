// +build wireinject

package copper

import (
	"github.com/google/wire"
	"github.com/tusharsoni/copper/v2/cconfig"
	"github.com/tusharsoni/copper/v2/clogger"
)

// InitApp creates a new Copper app along with its dependencies.
func InitApp() (*App, error) {
	panic(
		wire.Build(
			New,
			NewFlags,
			NewLifecycle,
			cconfig.New,
			clogger.NewWithConfig,

			wire.FieldsOf(new(*Flags), "Env", "ConfigDir"),
		),
	)
}

// WireModule can be used as part of google/wire setup to include the app's
// lifecycle, config, and logger.
var WireModule = wire.NewSet(
	wire.FieldsOf(new(*App), "Lifecycle", "Config", "Logger"),
)
