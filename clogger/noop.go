package clogger

// NewNoop returns a no-op implementation of Logger.
// Useful in passing it as a valid logger in unit tests.
func NewNoop() Logger {
	return &noop{}
}

type noop struct{}

func (l *noop) WithTags(tags map[string]interface{}) Logger {
	return l
}

func (l *noop) Debug(msg string) {}

func (l *noop) Info(msg string) {}

func (l *noop) Warn(msg string, err error) {}

func (l *noop) Error(msg string, err error) {}
