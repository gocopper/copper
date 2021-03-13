package cconfig_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tusharsoni/copper/cconfig"
	"github.com/tusharsoni/copper/cconfig/cconfigtest"
)

const (
	envTest = "test"
)

func TestNewConfig(t *testing.T) {
	t.Parallel()

	configDir := cconfigtest.SetupDirWithConfigs(t, "", "", "")

	config, err := cconfig.New(configDir, envTest)

	assert.NoError(t, err)
	assert.NotNil(t, config)
}

func TestConfig_Load_All(t *testing.T) {
	t.Parallel()

	var testConfig struct {
		Default string `default:"default"`
		Base    string `default:"default"`
		Env     string `default:"default"`
		Secrets string `default:"default"`
	}

	base := `
[group1]
base = "base"
`

	env := `
[group1]
env = "env"
`

	secrets := `
[group1]
secrets = "secrets"
`

	dir := cconfigtest.SetupDirWithConfigs(t, base, env, secrets)

	config, err := cconfig.New(dir, envTest)
	assert.NoError(t, err)

	err = config.Load("group1", &testConfig)
	assert.NoError(t, err)

	assert.Equal(t, "default", testConfig.Default)
	assert.Equal(t, "base", testConfig.Base)
	assert.Equal(t, "env", testConfig.Env)
	assert.Equal(t, "secrets", testConfig.Secrets)
}
