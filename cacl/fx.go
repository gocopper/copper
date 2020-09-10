package cacl

import "go.uber.org/fx"

// Fx module for the cacl package that provides the SQL implementation for all services.
var Fx = fx.Provide(
	newSQLRepo,
	newSvcImpl,
)
