package cmetrics

import "github.com/google/wire"

var WireModule = wire.NewSet(
	NewMetrics,
)
