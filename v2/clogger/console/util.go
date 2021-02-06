package console

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
