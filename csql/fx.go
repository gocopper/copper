package csql

import "go.uber.org/fx"

// Fx module that provides a connection to a Postgres DB using GORM.
var Fx = fx.Provide(
	newGormDB,
)
