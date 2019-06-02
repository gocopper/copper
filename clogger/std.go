package clogger

import (
	"log"

	"github.com/tusharsoni/copper/cerror"
)

type stdLogger struct{}

func newStdLogger() Logger {
	return &stdLogger{}
}

func (s *stdLogger) Debug(msg string, tags map[string]interface{}) {
	s.log(levelDebug, cerror.New(nil, msg, tags))
}

func (s *stdLogger) Info(msg string, tags map[string]interface{}) {
	s.log(levelInfo, cerror.New(nil, msg, tags))
}

func (s *stdLogger) Warn(msg string, err error) {
	s.log(levelWarn, cerror.New(err, msg, nil))
}

func (s *stdLogger) Error(msg string, err error) {
	s.log(levelError, cerror.New(err, msg, nil))
}

func (*stdLogger) log(lvl level, err error) {
	log.Printf("[%s] %s", string(lvl), err.Error())
}
