package clogger

import (
	"log"

	"go.uber.org/fx"

	"github.com/tusharsoni/copper/cerror"
)

type stdLoggerParams struct {
	fx.In

	Config Config `optional:"true"`
}

type stdLogger struct {
	config Config
}

func newStdLogger(p stdLoggerParams) Logger {
	if !p.Config.isValid() {
		p.Config = GetDefaultConfig()
	}

	return &stdLogger{
		config: p.Config,
	}
}

func (s *stdLogger) Debug(msg string, tags map[string]interface{}) {
	s.log(LevelDebug, cerror.New(nil, msg, tags))
}

func (s *stdLogger) Info(msg string, tags map[string]interface{}) {
	s.log(LevelInfo, cerror.New(nil, msg, tags))
}

func (s *stdLogger) Warn(msg string, err error) {
	s.log(LevelWarn, cerror.New(err, msg, nil))
}

func (s *stdLogger) Error(msg string, err error) {
	s.log(LevelError, cerror.New(err, msg, nil))
}

func (s *stdLogger) log(lvl Level, err error) {
	if lvl < s.config.MinLevel {
		return
	}

	log.Printf("[%s] %s", lvl.String(), err.Error())
}
