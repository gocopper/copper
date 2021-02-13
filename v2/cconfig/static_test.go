package cconfig_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tusharsoni/copper/v2/cconfig"
)

func TestNewStaticConfig(t *testing.T) {
	t.Parallel()

	config := cconfig.NewStaticConfig(nil)

	_, ok := config.(cconfig.Config)

	assert.NotNil(t, config)
	assert.True(t, ok)
}

func TestStaticConfig_Value(t *testing.T) {
	t.Parallel()

	config := cconfig.NewStaticConfig(map[string]interface{}{
		"key": "val",
	})

	val, ok := config.Value("key").(string)
	assert.True(t, ok)
	assert.Equal(t, "val", val)
}
