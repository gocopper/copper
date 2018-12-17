// Package chttp provides a fx module that can be used to create a http copper application.
package chttp

import "go.uber.org/fx"

// Fx provides the module for chttp that can be used to create a copper app.
// This module is provided by default when creating a http copper app.
var Fx = fx.Provide(
	newBodyReader,
	newResponder,

	newServer,
	newRouter,

	newHealthRoute,
)
