package csql

import "go.uber.org/fx"

// Fx module that provides a connection to a Postgres DB using GORM.
// It also provides a middleware to wrap http requests into db transactions.
var Fx = fx.Provide(
	newGormDB,

	newDBTxnMiddleware,
)
