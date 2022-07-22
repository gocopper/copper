package clogger

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

func tagsToKVs(tags map[string]interface{}) []interface{} {
	const TagsToKVsMultiplier = 2

	kvs := make([]interface{}, 0, len(tags)*TagsToKVsMultiplier)
	for k, v := range tags {
		kvs = append(kvs, k, v)
	}
	return kvs
}

func formatToZapEncoding(f Format) string {
	switch f {
	case FormatJSON:
		return "json"
	case FormatPlain:
		return "console"
	default:
		return "console"
	}
}
