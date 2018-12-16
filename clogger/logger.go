package clogger

const (
	levelDebug = level("DEBUG")
	levelInfo  = level("INFO")
	levelWarn  = level("WARN")
	levelError = level("ERROR")
)

type level string

// Logger can be used to log messages and errors.
type Logger interface {
	Debug(msg string, tags map[string]string)
	Info(msg string, tags map[string]string)
	Warn(msg string, err error)
	Error(msg string, err error)
}
