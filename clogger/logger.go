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

type CoreLogger interface {
	WithTags(tags map[string]any) Logger
	WithPrefix(prefix string) Logger

	Debug(msg string)
	Info(msg string)
	Warn(msg string, err error)
	Error(msg string, err error)
}

type Logger CoreLogger

func NewWithWriters(out, errWriter io.Writer, format Format, redactFields []string, levelFilter map[Level]bool, hooks []Hook) Logger {
	return &LoggerImpl{
		out:          out,
		err:          errWriter,
		tags:         make(map[string]any),
		format:       format,
		redactFields: expandRedactedFields(redactFields),
		levelFilter:  levelFilter,
		hooks:        hooks,
	}
}

func NewCore(config Config) (CoreLogger, error) {
	const LogFilePerms = 0666

	var (
		outFile io.Writer = os.Stdout
		errFile io.Writer = os.Stderr
		err     error
	)

	if config.Out != "" {
		outFile, err = os.OpenFile(config.Out, os.O_APPEND|os.O_CREATE|os.O_WRONLY, LogFilePerms)
		if err != nil {
			return nil, cerrors.New(err, "failed to open log file", map[string]any{
				"path": config.Out,
			})
		}
	}

	if config.Out == config.Err {
		errFile = outFile
	} else if config.Err != "" {
		errFile, err = os.OpenFile(config.Err, os.O_APPEND|os.O_CREATE|os.O_WRONLY, LogFilePerms)
		if err != nil {
			return nil, cerrors.New(err, "failed to open error log file", map[string]any{
				"path": config.Err,
			})
		}
	}

	var levelFilter map[Level]bool
	if len(config.LevelFilter) > 0 {
		levelFilter = make(map[Level]bool)
		for _, lvl := range config.LevelFilter {
			levelFilter[ParseLevel(lvl)] = true
		}
	}

	return &LoggerImpl{
		out:          outFile,
		err:          errFile,
		tags:         make(map[string]any),
		format:       config.Format,
		redactFields: expandRedactedFields(config.RedactFields),
		levelFilter:  levelFilter,
		hooks:        make([]Hook, 0),
	}, nil
}

func New(config Config, hooks []Hook) (Logger, error) {
	const LogFilePerms = 0666

	var (
		outFile io.Writer = os.Stdout
		errFile io.Writer = os.Stderr
		err     error
	)

	if config.Out != "" {
		outFile, err = os.OpenFile(config.Out, os.O_APPEND|os.O_CREATE|os.O_WRONLY, LogFilePerms)
		if err != nil {
			return nil, cerrors.New(err, "failed to open log file", map[string]any{
				"path": config.Out,
			})
		}
	}

	if config.Out == config.Err {
		errFile = outFile
	} else if config.Err != "" {
		errFile, err = os.OpenFile(config.Err, os.O_APPEND|os.O_CREATE|os.O_WRONLY, LogFilePerms)
		if err != nil {
			return nil, cerrors.New(err, "failed to open error log file", map[string]any{
				"path": config.Err,
			})
		}
	}

	var levelFilter map[Level]bool
	if len(config.LevelFilter) > 0 {
		levelFilter = make(map[Level]bool)
		for _, lvl := range config.LevelFilter {
			levelFilter[ParseLevel(lvl)] = true
		}
	}

	return &LoggerImpl{
		out:          outFile,
		err:          errFile,
		tags:         make(map[string]any),
		format:       config.Format,
		redactFields: expandRedactedFields(config.RedactFields),
		levelFilter:  levelFilter,
		hooks:        hooks,
	}, nil
}

type LoggerImpl struct {
	out          io.Writer
	err          io.Writer
	tags         map[string]any
	format       Format
	redactFields []string
	prefix       string
	levelFilter  map[Level]bool
	hooks        []Hook
}

func (l *LoggerImpl) WithTags(tags map[string]any) Logger {
	return &LoggerImpl{
		out:          l.out,
		err:          l.err,
		tags:         mergeTags(l.tags, tags),
		format:       l.format,
		redactFields: l.redactFields,
		prefix:       l.prefix,
		levelFilter:  l.levelFilter,
		hooks:        l.hooks,
	}
}

func (l *LoggerImpl) WithPrefix(prefix string) Logger {
	return &LoggerImpl{
		out:          l.out,
		err:          l.err,
		tags:         l.tags,
		format:       l.format,
		redactFields: l.redactFields,
		prefix:       prefix,
		levelFilter:  l.levelFilter,
		hooks:        l.hooks,
	}
}

func (l *LoggerImpl) Debug(msg string) {
	l.log(l.out, LevelDebug, msg, nil) //nolint:goerr113
}

func (l *LoggerImpl) Info(msg string) {
	l.log(l.out, LevelInfo, msg, nil) //nolint:goerr113
}

func (l *LoggerImpl) Warn(msg string, err error) {
	l.log(l.err, LevelWarn, msg, err)
}

func (l *LoggerImpl) Error(msg string, err error) {
	l.log(l.err, LevelError, msg, err)
}

func (l *LoggerImpl) log(dest io.Writer, lvl Level, msg string, err error) {
	// Filter by level (if level filter is set)
	if l.levelFilter != nil && !l.levelFilter[lvl] {
		return
	}

	switch l.format {
	case FormatJSON:
		l.logJSON(dest, lvl, msg, err)
	case FormatPlain:
		fallthrough
	default:
		l.logPlain(dest, lvl, msg, err)
	}

	for i := range l.hooks {
		l.hooks[i].OnLog(lvl, msg, l.tags, err)
	}
}

func (l *LoggerImpl) logJSON(dest io.Writer, lvl Level, msg string, err error) {
	if l.prefix != "" {
		msg = "[" + l.prefix + "] " + msg
	}

	var dict = map[string]any{
		"ts":    time.Now().Format(time.RFC3339),
		"level": lvl.String(),
		"msg":   msg,
	}

	if err != nil {
		dict["error"] = cerrors.WithoutTags(err).Error()
	}

	if redactedTags, err := redactJSONObject(mergeTags(cerrors.Tags(err), l.tags), l.redactFields); err != nil {
		dict["tags"] = cerrors.New(err, "tag redaction failed", nil).Error()
	} else {
		dict["tags"] = redactedTags
	}

	enc := json.NewEncoder(dest)
	enc.SetEscapeHTML(false)

	_ = enc.Encode(dict)
}

func (l *LoggerImpl) logPlain(dest io.Writer, lvl Level, msg string, err error) {
	if l.prefix != "" {
		msg = "[" + l.prefix + "] " + msg
	}

	var (
		logErr = cerrors.New(nil, msg, l.tags).Error()

		o strings.Builder
	)

	if len(l.redactFields) == 0 {
		o.WriteString(logErr)

		if err != nil {
			o.WriteString(" because\n> ")
			o.WriteString(err.Error())
		}
	} else {
		o.WriteString("<field redacting not supported for plain logs>")
	}

	log.New(dest, "", log.LstdFlags).Printf("[%s] %s", lvl.String(), o.String())
}
