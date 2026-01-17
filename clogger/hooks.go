package clogger

type Hook interface {
	OnLog(level Level, msg string, tags map[string]any, err error)
}
