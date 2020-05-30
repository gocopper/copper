package cauth2

import "go.uber.org/fx"

var Fx = fx.Provide(
	NewSQLRepo,
	NewSvc,
)

var FxMigrations = fx.Invoke(
	RunMigrations,
)
