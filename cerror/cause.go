package cerror

func Cause(err error) error {
	if cerr, ok := err.(Error); ok {
		return cerr.Cause
	}
	return nil
}
