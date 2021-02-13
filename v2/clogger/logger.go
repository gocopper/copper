package clogger

import (
	"errors"
	"io"
	"log"
	"os"

	"github.com/tusharsoni/copper/cerror"
	"github.com/tusharsoni/copper/v2/cconfig"
	"github.com/tusharsoni/copper/v2/cerrors"
)

// Logger can be used to log messages and errors.
type Logger interface {
	WithTags(tags map[string]interface{}) Logger

	Debug(msg string)
	Info(msg string)
	Warn(msg string, err error)
	Error(msg string, err error)
}

// New returns a Logger implementation that can logs to console.
func New() Logger {
	return NewWithWriters(os.Stdout, os.Stderr)
}

// NewWithConfig creates a Logger based on the provided config.
// Example TOML config:
// [clogger]
// out = "./logs.out"
// err = "./logs.err".
func NewWithConfig(config cconfig.Config) (Logger, error) {
	var (
		outFile io.Writer = os.Stdout
		errFile io.Writer = os.Stderr
		err     error
	)

	outFilePath, ok := config.Value("clogger.out").(string)

	if ok {
		outFile, err = os.OpenFile(outFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666) //nolint:gosec
		if err != nil {
			return nil, cerrors.New(err, "failed to open log file", map[string]interface{}{
				"path": outFilePath,
			})
		}
	}

	errFilePath, ok := config.Value("clogger.err").(string)

	if ok && errFilePath == outFilePath {
		errFile = outFile
	} else if ok {
		errFile, err = os.OpenFile(errFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666) //nolint:gosec
		if err != nil {
			return nil, cerrors.New(err, "failed to open error log file", map[string]interface{}{
				"path": errFilePath,
			})
		}
	}

	return NewWithWriters(outFile, errFile), nil
}

// NewWithWriters creates a Logger that uses the provided writers. out is
// used for debug and info levels. err is used for warn and error levels.
func NewWithWriters(out, err io.Writer) Logger {
	return &logger{
		out:  log.New(out, "", log.LstdFlags),
		err:  log.New(err, "", log.LstdFlags),
		tags: make(map[string]interface{}),
	}
}

type logger struct {
	out  *log.Logger
	err  *log.Logger
	tags map[string]interface{}
}

func (l *logger) WithTags(tags map[string]interface{}) Logger {
	return &logger{
		out:  l.out,
		err:  l.err,
		tags: mergeTags(l.tags, tags),
	}
}

func (l *logger) Debug(msg string) {
	l.log(l.out, LevelDebug, errors.New(msg)) //nolint:goerr113
}

func (l *logger) Info(msg string) {
	l.log(l.out, LevelInfo, errors.New(msg)) //nolint:goerr113
}

func (l *logger) Warn(msg string, err error) {
	l.log(l.err, LevelWarn, cerror.New(err, msg, nil))
}

func (l *logger) Error(msg string, err error) {
	l.log(l.err, LevelError, cerror.New(err, msg, nil))
}

func (l *logger) log(logger *log.Logger, lvl Level, err error) {
	logger.Printf("[%s] %s", lvl.String(), cerror.WithTags(err, l.tags).Error())
}
