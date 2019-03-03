package cacl

import "go.uber.org/fx"

// Fx module for the cacl package that provides the SQL implementation for all services.
var Fx = fx.Provide(
	newSQLRepo,
	newSvcImpl,
)

// RunMigrations can be used with fx.Invoke to run the db migrations for SQL implementation of cacl
var RunMigrations = fx.Invoke(
	runMigrations,
)
