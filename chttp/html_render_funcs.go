package chttp

import "errors"

// https://github.com/gohugoio/hugo/blob/a2670bf460e10ed5de69f90abbe7c4e2b32068cf/tpl/collections/collections.go#L149
func dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dictionary call")
	}

	root := make(map[string]interface{})

	for i := 0; i < len(values); i += 2 {
		dict := root
		var key string
		switch v := values[i].(type) {
		case string:
			key = v
		case []string:
			for i := 0; i < len(v)-1; i++ {
				key = v[i]
				var m map[string]interface{}
				v, found := dict[key]
				if found {
					m = v.(map[string]interface{})
				} else {
					m = make(map[string]interface{})
					dict[key] = m
				}
				dict = m
			}
			key = v[len(v)-1]
		default:
			return nil, errors.New("invalid dictionary key")
		}
		dict[key] = values[i+1]
	}

	return root, nil
}
