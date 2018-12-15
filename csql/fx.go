package csql

import "go.uber.org/fx"

var Fx = fx.Provide(
	NewGormDB,
)
