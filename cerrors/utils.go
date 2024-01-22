package cerrors

import "errors"

func Tags(err error) map[string]interface{} {
	var cerr Error

	if !errors.As(err, &cerr) {
		return nil
	}

	return mergeTags(cerr.Tags, Tags(cerr.Cause))
}

func mergeTags(t1, t2 map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})

	for k, v := range t1 {
		merged[k] = v
	}

	for k, v := range t2 {
		merged[k] = v
	}

	return merged
}
