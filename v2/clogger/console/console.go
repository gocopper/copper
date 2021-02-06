// Package console provides an implementation to clogger.Logger that logs to the console with
// human-friendly formatting
package console

import (
	"errors"
	"log"

	"github.com/tusharsoni/copper/cerror"
	"github.com/tusharsoni/copper/v2/clogger"
)

// New returns a clogger.Logger implementation that is best suited to log to console.
func New() clogger.Logger {
	return &logger{
		tags: make(map[string]interface{}),
	}
}

type logger struct {
	tags map[string]interface{}
}

func (l *logger) WithTags(tags map[string]interface{}) clogger.Logger {
	return &logger{
		tags: mergeTags(l.tags, tags),
	}
}

func (l *logger) Debug(msg string) {
	l.log(clogger.LevelDebug, errors.New(msg)) //nolint:goerr113
}

func (l *logger) Info(msg string) {
	l.log(clogger.LevelInfo, errors.New(msg)) //nolint:goerr113
}

func (l *logger) Warn(msg string, err error) {
	l.log(clogger.LevelWarn, cerror.New(err, msg, nil))
}

func (l *logger) Error(msg string, err error) {
	l.log(clogger.LevelError, cerror.New(err, msg, nil))
}

func (l *logger) log(lvl clogger.Level, err error) {
	log.Printf("[%s] %s", lvl.String(), cerror.WithTags(err, l.tags).Error())
}