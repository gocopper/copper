package clogger

// NewRecorder returns an implementation of Logger that keeps
// a record of each log. Useful in unit tests when logs need
// to be tested.
func NewRecorder(logs *[]RecordedLog) Logger {
	return &recorder{
		Logs: logs,
		tags: make(map[string]any),
	}
}

// RecordedLog represents a single log.
type RecordedLog struct {
	Level  Level
	Tags   map[string]any
	Msg    string
	Error  error
	Prefix string
}

type recorder struct {
	Logs   *[]RecordedLog
	tags   map[string]any
	prefix string
}

func (l *recorder) WithTags(tags map[string]any) Logger {
	return &recorder{
		Logs:   l.Logs,
		tags:   mergeTags(l.tags, tags),
		prefix: l.prefix,
	}
}

func (l *recorder) WithPrefix(prefix string) Logger {
	return &recorder{
		Logs:   l.Logs,
		tags:   l.tags,
		prefix: prefix,
	}
}

func (l *recorder) Debug(msg string) {
	*l.Logs = append(*l.Logs, RecordedLog{
		Level:  LevelDebug,
		Tags:   mergeTags(l.tags, nil),
		Msg:    msg,
		Error:  nil,
		Prefix: l.prefix,
	})
}

func (l *recorder) Info(msg string) {
	*l.Logs = append(*l.Logs, RecordedLog{
		Level:  LevelInfo,
		Tags:   mergeTags(l.tags, nil),
		Msg:    msg,
		Error:  nil,
		Prefix: l.prefix,
	})
}

func (l *recorder) Warn(msg string, err error) {
	*l.Logs = append(*l.Logs, RecordedLog{
		Level:  LevelWarn,
		Tags:   mergeTags(l.tags, nil),
		Msg:    msg,
		Error:  err,
		Prefix: l.prefix,
	})
}

func (l *recorder) Error(msg string, err error) {
	*l.Logs = append(*l.Logs, RecordedLog{
		Level:  LevelError,
		Tags:   mergeTags(l.tags, nil),
		Msg:    msg,
		Error:  err,
		Prefix: l.prefix,
	})
}
