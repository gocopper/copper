package clogger

import (
	"bytes"
	"encoding/json"
	"strings"
)

func redactJSONObject(in map[string]any, redactFields []string) (map[string]any, error) {
	var b bytes.Buffer

	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)

	err := enc.Encode(in)
	if err != nil {
		return nil, err
	}

	redacted, err := redactJSON(b.Bytes(), redactFields)
	if err != nil {
		return nil, err
	}

	var out map[string]any
	err = json.Unmarshal(redacted, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func redactJSON(in json.RawMessage, redactFields []string) (json.RawMessage, error) {
	var err error

	if in[0] == 123 { //  123 is `{` => object
		var cont map[string]json.RawMessage

		err = json.Unmarshal(in, &cont)
		if err != nil {
			return nil, err
		}

		for k, v := range cont {
			didRedact := false
			for i := range redactFields {
				if strings.Contains(strings.ToLower(k), strings.ToLower(redactFields[i])) {
					cont[k] = json.RawMessage(`"redacted"`)
					didRedact = true
					break
				}
			}

			if didRedact {
				continue
			}

			cont[k], err = redactJSON(v, redactFields)
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
			cont[i], err = redactJSON(v, redactFields)
			if err != nil {
				return nil, err
			}
		}

		return json.Marshal(cont)
	}

	return in, nil
}
