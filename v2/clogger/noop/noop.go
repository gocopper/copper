// Package noop provides an implementation to clogger.Logger that does nothing.
// Useful in passing it as a valid logger in unit tests.
package noop

import "github.com/tusharsoni/copper/v2/clogger"

// New returns a no-op implementation of clogger.Logger.
func New() clogger.Logger {
	return &noop{}
}

type noop struct{}

func (l *noop) WithTags(tags map[string]interface{}) clogger.Logger {
	return l
}

func (l *noop) Debug(msg string) {}

func (l *noop) Info(msg string) {}

func (l *noop) Warn(msg string, err error) {}

func (l *noop) Error(msg string, err error) {}
