package cauth

import (
	"github.com/google/wire"
)

// WireModule can be used as part of google/wire setup.
var WireModule = wire.NewSet( // nolint:gochecknoglobals
	NewSvc,
	NewRepo,
	NewMigration,
	NewSetSessionMiddleware,
	NewVerifySessionMiddleware,

	wire.Struct(new(NewRouterParams), "*"),
	NewRouter,
)
