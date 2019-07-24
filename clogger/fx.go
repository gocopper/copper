// Package clogger provides an interface to log messages.
package clogger

import "go.uber.org/fx"

// StdFx provides a logger implementation using the stdlib log package
var StdFx = fx.Provide(
	newStdLogger,
)

// SentryFx provides a logger implementation that sends messages to Sentry
var SentryFx = fx.Provide(
	newSentryLogger,
)
