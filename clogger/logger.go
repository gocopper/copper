package clogger

const (
	DEBUG = Level("DEBUG")
	INFO  = Level("INFO")
	WARN  = Level("WARN")
	ERROR = Level("ERROR")
)

type Level string

type Logger interface {
	Debug(msg string, tags map[string]string)
	Info(msg string, tags map[string]string)
	Warn(msg string, err error)
	Error(msg string, err error)
}
