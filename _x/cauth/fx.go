package cauth

import "go.uber.org/fx"

var Fx = fx.Provide(
	NewSQLRepo,
	NewSvc,

	NewAuthMiddleware,
)

var FxMigrations = fx.Invoke(
	RunMigrations,
)
