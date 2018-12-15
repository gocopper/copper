package chttp

import "go.uber.org/fx"

var Fx = fx.Provide(
	NewServer,
	NewRouter,
	NewHealthRoute,
)
