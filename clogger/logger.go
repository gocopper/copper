package clogger

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"
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
	return NewWithWriters(os.Stdout, os.Stderr, FormatPlain, nil)
}

// NewWithConfig creates a Logger based on the provided config.
func NewWithConfig(config Config) (Logger, error) {
	const LogFilePerms = 0666

	var (
		outFile io.Writer = os.Stdout
		errFile io.Writer = os.Stderr
		err     error
	)

	if config.Out != "" {
		outFile, err = os.OpenFile(config.Out, os.O_APPEND|os.O_CREATE|os.O_WRONLY, LogFilePerms)
		if err != nil {
			return nil, cerrors.New(err, "failed to open log file", map[string]interface{}{
				"path": config.Out,
			})
		}
	}

	if config.Out == config.Err {
		errFile = outFile
	} else if config.Err != "" {
		errFile, err = os.OpenFile(config.Err, os.O_APPEND|os.O_CREATE|os.O_WRONLY, LogFilePerms)
		if err != nil {
			return nil, cerrors.New(err, "failed to open error log file", map[string]interface{}{
				"path": config.Err,
			})
		}
	}

	return NewWithWriters(outFile, errFile, config.Format, config.RedactFields), nil
}

// NewWithWriters creates a Logger that uses the provided writers. out is
// used for debug and info levels. err is used for warn and error levels.
func NewWithWriters(out, err io.Writer, format Format, redactFields []string) Logger {
	return &logger{
		out:          out,
		err:          err,
		tags:         make(map[string]interface{}),
		format:       format,
		redactFields: expandRedactedFields(redactFields),
	}
}

type logger struct {
	out          io.Writer
	err          io.Writer
	tags         map[string]interface{}
	format       Format
	redactFields []string
}

func (l *logger) WithTags(tags map[string]interface{}) Logger {
	return &logger{
		out:          l.out,
		err:          l.err,
		tags:         mergeTags(l.tags, tags),
		format:       l.format,
		redactFields: l.redactFields,
	}
}

func (l *logger) Debug(msg string) {
	l.log(l.out, LevelDebug, msg, nil) //nolint:goerr113
}

func (l *logger) Info(msg string) {
	l.log(l.out, LevelInfo, msg, nil) //nolint:goerr113
}

func (l *logger) Warn(msg string, err error) {
	l.log(l.err, LevelWarn, msg, err)
}

func (l *logger) Error(msg string, err error) {
	l.log(l.err, LevelError, msg, err)
}

func (l *logger) log(dest io.Writer, lvl Level, msg string, err error) {
	switch l.format {
	case FormatJSON:
		l.logJSON(dest, lvl, msg, err)
	case FormatPlain:
		fallthrough
	default:
		l.logPlain(dest, lvl, msg, err)
	}
}

func (l *logger) logJSON(dest io.Writer, lvl Level, msg string, err error) {
	var dict = map[string]interface{}{
		"ts":    time.Now().Format(time.RFC3339),
		"level": lvl.String(),
		"msg":   msg,
	}

	dict = mergeTags(mergeTags(dict, redactTags(l.tags, l.redactFields)), redactTags(cerrors.Tags(err), l.redactFields))

	if err != nil {
		errStr := err.Error()
		if stringHasRedactedFields(errStr, l.redactFields) {
			dict["error"] = "<redacted>"
		} else {
			dict["error"] = errStr
		}
	}

	enc := json.NewEncoder(dest)
	enc.SetEscapeHTML(false)

	_ = enc.Encode(dict)
}

func (l *logger) logPlain(dest io.Writer, lvl Level, msg string, err error) {
	var (
		logErr = cerrors.New(nil, msg, redactTags(l.tags, l.redactFields)).Error()

		o strings.Builder
	)

	o.WriteString(logErr)

	if err != nil {
		o.WriteString(" because\n> ")
		errStr := err.Error()
		if stringHasRedactedFields(errStr, l.redactFields) {
			o.WriteString("<redacted>")
		} else {
			o.WriteString(errStr)
		}
	}

	log.New(dest, "", log.LstdFlags).Printf("[%s] %s", lvl.String(), o.String())
}
