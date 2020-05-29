package ctexter

import "go.uber.org/fx"

var FxLogger = fx.Provide(
	newLoggerSvc,
)

var FxAWS = fx.Provide(
	newAWSSvc,
)
