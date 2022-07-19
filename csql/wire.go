package csql

import "github.com/google/wire"

// WireModule can be used as part of google/wire setup.
var WireModule = wire.NewSet(
	NewDBConnection,
	NewQuerier,
	NewMigrator,
	LoadConfig,
	NewTxMiddleware,

	wire.Struct(new(NewMigratorParams), "*"),
)
