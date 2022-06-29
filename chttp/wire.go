package chttp

import (
	"github.com/google/wire"
)

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

// WireModuleEmptyHTML provides empty/default values for html and static dirs. This can be used to satisfy
// wire when the project does not use/need html rendering.
var WireModuleEmptyHTML = wire.NewSet( //nolint:gochecknoglobals
	wire.InterfaceValue(new(HTMLDir), &EmptyFS{}),
	wire.InterfaceValue(new(StaticDir), &EmptyFS{}),
	wire.Value([]HTMLRenderFunc{}),
)
