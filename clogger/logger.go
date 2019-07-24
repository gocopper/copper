package clogger

// Logger can be used to log messages and errors.
type Logger interface {
	WithTags(tags map[string]interface{}) Logger

	Debug(msg string)
	Info(msg string)
	Warn(msg string, err error)
	Error(msg string, err error)
}
