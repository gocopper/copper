package cconfig_test

import (
	"path"
	"testing"

	"github.com/gocopper/copper/cconfig"
	"github.com/gocopper/copper/cconfig/cconfigtest"
	"github.com/stretchr/testify/assert"
)

func TestLoader_Load(t *testing.T) {
	t.Parallel()

	var (
		dir = cconfigtest.SetupDirWithConfigs(t, map[string]string{
			"base.toml": `
				[group1]
				key1 = "val1"
			`,
			"test.toml": `
				extends = "base.toml"

				[group1]
				key1 = "val2"
			`,
		})
		fp = cconfig.Path(path.Join(dir, "test.toml"))
	)

	t.Run("load with extends and key overrides", func(t *testing.T) {
		t.Parallel()

		var testConfig struct {
			Key1 string `toml:"key1"`
		}

		configs, err := cconfig.NewWithKeyOverrides(fp)
		assert.NoError(t, err)

		err = configs.Load("group1", &testConfig)
		assert.NoError(t, err)

		assert.Equal(t, "val2", testConfig.Key1)
	})

	t.Run("error with key overrides disabled", func(t *testing.T) {
		t.Parallel()

		_, err := cconfig.New(fp)

		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "key is being overridden when key overrides are disabled")
	})
}
