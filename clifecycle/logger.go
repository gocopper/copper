package clifecycle

// Logger provides the methods needed by Lifecycle to log errors.
type Logger interface {
	Error(msg string, err error)
}
