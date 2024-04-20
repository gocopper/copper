package clogger

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRedactJSON(t *testing.T) {
	var t1 = map[string]interface{}{
		"a": 1,
		"b": "foo",
		"c": map[string]interface{}{"d": 2},
		"e": []interface{}{1, 2, map[string]interface{}{"f": 3}},
	}

	in, err := json.Marshal(t1)
	assert.NoError(t, err)

	out, err := redactJSON(in, map[string]bool{
		"f": true,
	})
	assert.NoError(t, err)

	assert.Equal(t, `{"a":1,"b":"foo","c":{"d":2},"e":[1,2,{"f":"redacted"}]}`, string(out))
}
