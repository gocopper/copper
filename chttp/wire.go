package chttp

import "github.com/google/wire"

// WireModule can be used as part of google/wire setup.
var WireModule = wire.NewSet( //nolint:gochecknoglobals
	NewReaderWriter,
	NewRequestLoggerMiddleware,
	wire.Struct(new(NewServerParams), "*"),
	NewServer,
	NewStaticRouter,
)
