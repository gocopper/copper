package anonymous

import "go.uber.org/fx"

var Fx = fx.Provide(
	NewSvc,

	NewRouter,
	NewCreateSessionRoute,
)
