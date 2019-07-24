package clogger

// Logger can be used to log messages and errors.
type Logger interface {
	Debug(msg string, tags map[string]interface{})
	Info(msg string, tags map[string]interface{})
	Warn(msg string, err error)
	Error(msg string, err error)
}
