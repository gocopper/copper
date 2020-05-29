package clogger

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/tusharsoni/copper/cerror"
	"go.uber.org/fx"
)

type sentryLoggerParams struct {
	fx.In

	LC     fx.Lifecycle
	Config SentryConfig `optional:"true"`
}

type sentryLogger struct {
	tags   map[string]interface{}
	logger *stdLogger
	config SentryConfig
}

func newSentryLogger(p sentryLoggerParams) (Logger, error) {
	if !p.Config.isValid() {
		p.Config = GetDefaultSentryConfig()
	}

	err := sentry.Init(sentry.ClientOptions{Dsn: p.Config.Dsn})
	if err != nil {
		return nil, cerror.New(err, "failed to init sentry", nil)
	}

	p.LC.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			didFlush := sentry.Flush(time.Second * 5)
			if !didFlush {
				return errors.New("failed to flush sentry logs")
			}
			return nil
		},
	})

	return &sentryLogger{
		tags: make(map[string]interface{}),
		logger: &stdLogger{
			config: StdConfig{
				MinLevel: p.Config.MinLevelForStd,
			},
		},
		config: p.Config,
	}, nil
}

func (s *sentryLogger) WithTags(tags map[string]interface{}) Logger {
	return &sentryLogger{
		tags:   mergeTags(s.tags, tags),
		logger: s.logger,
		config: s.config,
	}
}

func (s *sentryLogger) Debug(msg string) {
	s.log(LevelDebug, msg, nil)
}

func (s *sentryLogger) Info(msg string) {
	s.log(LevelInfo, msg, nil)
}

func (s *sentryLogger) Warn(msg string, err error) {
	s.log(LevelWarn, msg, err)
}

func (s *sentryLogger) Error(msg string, err error) {
	s.log(LevelError, msg, err)
}

func (s *sentryLogger) log(level Level, msg string, err error) {
	if level >= s.config.MinLevelForStd {
		switch level {
		case LevelDebug:
			s.logger.WithTags(s.tags).Debug(msg)
		case LevelInfo:
			s.logger.WithTags(s.tags).Info(msg)
		case LevelWarn:
			s.logger.WithTags(s.tags).Warn(msg, err)
		case LevelError:
			s.logger.WithTags(s.tags).Error(msg, err)
		}
	}

	if level < s.config.MinLevelForCapture || s.isMsgIgnored(msg) {
		return
	}

	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetLevel(level.toSentryLevel())
		for k, v := range s.tags {
			rv := reflect.ValueOf(v)
			if rv.Kind() == reflect.Ptr && !rv.IsNil() {
				scope.SetTag(k, fmt.Sprintf("%+v", rv.Elem()))
			} else {
				scope.SetTag(k, fmt.Sprintf("%+v", v))
			}
		}
		if err != nil {
			scope.SetExtra("_error", err.Error())
		}
		sentry.CaptureMessage(msg)
	})
}

func (s *sentryLogger) isMsgIgnored(msg string) bool {
	for _, ignoredMsg := range s.config.IgnoredMsgsForCapture {
		if msg == ignoredMsg {
			return true
		}
	}

	return false
}

func (l Level) toSentryLevel() sentry.Level {
	switch l {
	case LevelDebug:
		return sentry.LevelDebug
	case LevelInfo:
		return sentry.LevelInfo
	case LevelWarn:
		return sentry.LevelWarning
	case LevelError:
		return sentry.LevelError
	default:
		return sentry.LevelError
	}
}
