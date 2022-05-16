package chttp

import "github.com/google/wire"

// WireModule can be used as part of google/wire setup.
var WireModule = wire.NewSet( //nolint:gochecknoglobals
	LoadConfig,
	NewReaderWriter,
	NewRequestLoggerMiddleware,
	wire.Struct(new(NewServerParams), "*"),
	NewServer,
	wire.Struct(new(NewHTMLRouterParams), "*"),
	NewHTMLRouter,
	wire.Struct(new(NewHTMLRendererParams), "*"),
	NewHTMLRenderer,
)
