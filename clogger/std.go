package clogger

import (
	"log"

	"github.com/tusharsoni/copper/cerror"
)

type StdLogger struct{}

func NewStdLogger() Logger {
	return &StdLogger{}
}

func (s *StdLogger) Debug(msg string, tags map[string]string) {
	s.log(DEBUG, cerror.New(nil, msg, tags))
}

func (s *StdLogger) Info(msg string, tags map[string]string) {
	s.log(DEBUG, cerror.New(nil, msg, tags))
}

func (s *StdLogger) Warn(msg string, err error) {
	s.log(WARN, cerror.New(err, msg, nil))
}

func (s *StdLogger) Error(msg string, err error) {
	s.log(ERROR, cerror.New(err, msg, nil))
}

func (*StdLogger) log(lvl Level, err error) {
	log.Printf("[%s] %s", string(lvl), err.Error())
}
