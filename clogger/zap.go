package clogger

import (
	"context"

	"github.com/gocopper/copper/cerrors"
	"github.com/gocopper/copper/clifecycle"
	"go.uber.org/zap"
)

// NewZapLogger creates a Logger that internally uses go.uber.org/zap for logging
func NewZapLogger(config Config, lc *clifecycle.Lifecycle) (Logger, error) {
	const OutStdErr = "stderr"

	var (
		outPath    = OutStdErr
		errOutPath = OutStdErr
	)

	if config.Out != "" {
		outPath = config.Out
	}

	if config.Err != "" {
		errOutPath = config.Err
	}

	encoderConfig := zap.NewDevelopmentEncoderConfig()
	if config.Format == FormatJSON {
		encoderConfig = zap.NewProductionEncoderConfig()
	}

	z, err := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		Encoding:         formatToZapEncoding(config.Format),
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{outPath},
		ErrorOutputPaths: []string{errOutPath},
	}.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, cerrors.New(err, "failed to create zap logger", nil)
	}

	lc.OnStop(func(ctx context.Context) error {
		// Skip sync if logs are written to stderr because it will throw an error:
		// https://github.com/uber-go/zap/issues/880
		if outPath == OutStdErr && errOutPath == OutStdErr {
			return nil
		}

		return z.Sync()
	})

	return &zapLogger{
		zap:  z.Sugar(),
		tags: make(map[string]interface{}),
	}, nil
}

type zapLogger struct {
	zap  *zap.SugaredLogger
	tags map[string]interface{}
}

func (l *zapLogger) WithTags(tags map[string]interface{}) Logger {
	return &zapLogger{
		zap:  l.zap,
		tags: mergeTags(l.tags, tags),
	}
}

func (l *zapLogger) Debug(msg string) {
	l.zap.Debugw(msg, tagsToKVs(l.tags)...)
}

func (l *zapLogger) Info(msg string) {
	l.zap.Infow(msg, tagsToKVs(l.tags)...)
}

func (l *zapLogger) Warn(msg string, err error) {
	l.zap.With("error", err).Warnw(msg, tagsToKVs(l.tags)...)
}

func (l *zapLogger) Error(msg string, err error) {
	l.zap.With("error", err).Errorw(msg, tagsToKVs(l.tags)...)
}
