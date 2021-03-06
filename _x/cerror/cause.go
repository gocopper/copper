package cerror

// Cause returns the cause of the error as provided when the error was created using cerror.New.
// Cause returns nil if the error has no given cause.
func Cause(err error) error {
	if cerr, ok := err.(Error); ok {
		return cerr.Cause
	}
	return nil
}

// HasCause returns whether the target error is in the error chain
func HasCause(err, target error) bool {
	cerr, ok := err.(Error)
	if !ok {
		return err.Error() == target.Error()
	}

	if cerr.Cause == nil {
		return false
	}

	return HasCause(cerr.Cause, target)
}
