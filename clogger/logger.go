package clogger

const (
	DEBUG = Level("DEBUG")
	INFO  = Level("INFO")
	WARN  = Level("WARN")
	ERROR = Level("ERROR")
)

type Level string

type Logger interface {
	Debugw(msg string, tags map[string]string)
	Infow(msg string, tags map[string]string)
	Warnw(msg string, tags map[string]string)
	Errorw(msg string, tags map[string]string)
}
