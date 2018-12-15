package clogger

import "go.uber.org/fx"

var StdFx = fx.Provide(
	NewStdLogger,
)
