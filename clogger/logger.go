package clogger

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"time"

	"github.com/gocopper/copper/cerrors"
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
	return NewWithWriters(os.Stdout, os.Stderr, FormatPlain)
}

// NewWithConfig creates a Logger based on the provided config.
func NewWithConfig(config Config) (Logger, error) {
	var (
		outFile io.Writer = os.Stdout
		errFile io.Writer = os.Stderr
		err     error
	)

	if config.Out != "" {
		outFile, err = os.OpenFile(config.Out, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666) //nolint:gosec
		if err != nil {
			return nil, cerrors.New(err, "failed to open log file", map[string]interface{}{
				"path": config.Out,
			})
		}
	}

	if config.Out == config.Err {
		errFile = outFile
	} else if config.Err != "" {
		errFile, err = os.OpenFile(config.Err, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666) //nolint:gosec
		if err != nil {
			return nil, cerrors.New(err, "failed to open error log file", map[string]interface{}{
				"path": config.Err,
			})
		}
	}

	return NewWithWriters(outFile, errFile, config.Format), nil
}

// NewWithWriters creates a Logger that uses the provided writers. out is
// used for debug and info levels. err is used for warn and error levels.
func NewWithWriters(out, err io.Writer, format Format) Logger {
	return &logger{
		out:    out,
		err:    err,
		tags:   make(map[string]interface{}),
		format: format,
	}
}

type logger struct {
	out    io.Writer
	err    io.Writer
	tags   map[string]interface{}
	format Format
}

func (l *logger) WithTags(tags map[string]interface{}) Logger {
	return &logger{
		out:    l.out,
		err:    l.err,
		tags:   mergeTags(l.tags, tags),
		format: l.format,
	}
}

func (l *logger) Debug(msg string) {
	l.log(l.out, LevelDebug, errors.New(msg)) //nolint:goerr113
}

func (l *logger) Info(msg string) {
	l.log(l.out, LevelInfo, errors.New(msg)) //nolint:goerr113
}

func (l *logger) Warn(msg string, err error) {
	l.log(l.err, LevelWarn, cerrors.New(err, msg, nil))
}

func (l *logger) Error(msg string, err error) {
	l.log(l.err, LevelError, cerrors.New(err, msg, nil))
}

func (l *logger) log(dest io.Writer, lvl Level, err error) {
	switch l.format {
	case FormatJSON:
		l.logJSON(dest, lvl, err)
	case FormatPlain:
		fallthrough
	default:
		l.logPlain(dest, lvl, err)
	}
}

func (l *logger) logJSON(dest io.Writer, lvl Level, err error) {
	var dict = map[string]interface{}{
		"ts":    time.Now().Format(time.RFC3339),
		"level": lvl.String(),
	}

	dict = mergeTags(dict, l.tags)

	switch cerr := err.(type) {
	case cerrors.Error:
		dict["msg"] = cerr.Message

		dict = mergeTags(dict, cerr.Tags)

		if cerr.Cause != nil {
			dict["error"] = cerr.Cause.Error()
		}
	default:
		dict["msg"] = err.Error()
	}

	jsonStr, _ := json.Marshal(dict)
	_, _ = dest.Write([]byte(string(jsonStr) + "\n"))
}

func (l *logger) logPlain(dest io.Writer, lvl Level, err error) {
	log.New(dest, "", log.LstdFlags).Printf("[%s] %s", lvl.String(), cerrors.WithTags(err, l.tags).Error())
}
