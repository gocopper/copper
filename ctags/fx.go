package ctags

import "go.uber.org/fx"

// Fx provides the ctags module that includes
// - a SQL repo to manage the storage of tags
// - a service implementation to manage tags
var Fx = fx.Provide(
	newSQLRepo,
	newSvcImpl,
)
