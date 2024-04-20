package clogger

import (
	"bytes"
	"encoding/json"
	"strings"
)

func redactJSONObject(in map[string]interface{}, redactFields []string) (map[string]interface{}, error) {
	var b bytes.Buffer

	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)

	err := enc.Encode(in)
	if err != nil {
		return nil, err
	}

	redactFieldsSet := make(map[string]bool)
	for _, f := range redactFields {
		redactFieldsSet[strings.ToLower(f)] = true
	}

	redacted, err := redactJSON(b.Bytes(), redactFieldsSet)
	if err != nil {
		return nil, err
	}

	var out map[string]interface{}
	err = json.Unmarshal(redacted, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func redactJSON(in json.RawMessage, redactKeys map[string]bool) (json.RawMessage, error) {
	var err error

	if in[0] == 123 { //  123 is `{` => object
		var cont map[string]json.RawMessage

		err = json.Unmarshal(in, &cont)
		if err != nil {
			return nil, err
		}

		for k, v := range cont {
			if redact, ok := redactKeys[strings.ToLower(k)]; ok && redact {
				cont[k] = json.RawMessage(`"redacted"`)
				continue
			}

			cont[k], err = redactJSON(v, redactKeys)
			if err != nil {
				return nil, err
			}
		}

		return json.Marshal(cont)
	} else if in[0] == 91 { // 91 is `[`  => array
		var cont []json.RawMessage

		err = json.Unmarshal(in, &cont)
		if err != nil {
			return nil, err
		}

		for i, v := range cont {
			cont[i], err = redactJSON(v, redactKeys)
			if err != nil {
				return nil, err
			}
		}

		return json.Marshal(cont)
	}

	return in, nil
}
