package clifecycle

// Logger provides the methods needed by Lifecycle to log errors.
type Logger interface {
	Info(msg string)
	Error(msg string, err error)
}
