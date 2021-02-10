package clogger

import (
	"errors"
	"log"

	"github.com/tusharsoni/copper/cerror"
)

// NewConsole returns a Logger implementation that is best suited to log to console
// with human-friendly formatting.
func NewConsole() Logger {
	return &console{
		tags: make(map[string]interface{}),
	}
}

type console struct {
	tags map[string]interface{}
}

func (l *console) WithTags(tags map[string]interface{}) Logger {
	return &console{
		tags: mergeTags(l.tags, tags),
	}
}

func (l *console) Debug(msg string) {
	l.log(LevelDebug, errors.New(msg)) //nolint:goerr113
}

func (l *console) Info(msg string) {
	l.log(LevelInfo, errors.New(msg)) //nolint:goerr113
}

func (l *console) Warn(msg string, err error) {
	l.log(LevelWarn, cerror.New(err, msg, nil))
}

func (l *console) Error(msg string, err error) {
	l.log(LevelError, cerror.New(err, msg, nil))
}

func (l *console) log(lvl Level, err error) {
	log.Printf("[%s] %s", lvl.String(), cerror.WithTags(err, l.tags).Error())
}
