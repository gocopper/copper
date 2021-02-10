package clogger

// NewRecorder returns an implementation of Logger that keeps
// a record of each log. Useful in unit tests when logs need
// to be tested.
func NewRecorder(logs *[]RecordedLog) Logger {
	return &recorder{
		Logs: logs,
		tags: make(map[string]interface{}),
	}
}

// RecordedLog represents a single log.
type RecordedLog struct {
	Level Level
	Tags  map[string]interface{}
	Msg   string
	Error error
}

type recorder struct {
	Logs *[]RecordedLog
	tags map[string]interface{}
}

func (l *recorder) WithTags(tags map[string]interface{}) Logger {
	return &recorder{
		Logs: l.Logs,
		tags: mergeTags(l.tags, tags),
	}
}

func (l *recorder) Debug(msg string) {
	*l.Logs = append(*l.Logs, RecordedLog{
		Level: LevelDebug,
		Tags:  mergeTags(l.tags, nil),
		Msg:   msg,
		Error: nil,
	})
}

func (l *recorder) Info(msg string) {
	*l.Logs = append(*l.Logs, RecordedLog{
		Level: LevelInfo,
		Tags:  mergeTags(l.tags, nil),
		Msg:   msg,
		Error: nil,
	})
}

func (l *recorder) Warn(msg string, err error) {
	*l.Logs = append(*l.Logs, RecordedLog{
		Level: LevelWarn,
		Tags:  mergeTags(l.tags, nil),
		Msg:   msg,
		Error: err,
	})
}

func (l *recorder) Error(msg string, err error) {
	*l.Logs = append(*l.Logs, RecordedLog{
		Level: LevelError,
		Tags:  mergeTags(l.tags, nil),
		Msg:   msg,
		Error: err,
	})
}
