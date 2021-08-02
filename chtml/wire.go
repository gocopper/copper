package chtml

import "github.com/google/wire"

// WireModule can be used as part of google/wire setup for fullstack apps
var WireModule = wire.NewSet( //nolint:gochecknoglobals
	wire.Struct(new(NewRouterParams), "*"),
	NewRouter,

	wire.Struct(new(NewRendererParams), "*"),
	NewRenderer,

	NewReaderWriter,
)
