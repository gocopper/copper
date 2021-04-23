package cconfig_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tusharsoni/copper/cconfig"
	"github.com/tusharsoni/copper/cconfig/cconfigtest"
)

const (
	envTest    = "test"
	projectDir = "/tmp/project-dir"
)

func TestNewConfig(t *testing.T) {
	t.Parallel()

	configDir := cconfigtest.SetupDirWithConfigs(t, "", "", "")

	config, err := cconfig.New(configDir, projectDir, envTest)

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

	config, err := cconfig.New(dir, projectDir, envTest)
	assert.NoError(t, err)

	err = config.Load("group1", &testConfig)
	assert.NoError(t, err)

	assert.Equal(t, "default", testConfig.Default)
	assert.Equal(t, "base", testConfig.Base)
	assert.Equal(t, "env", testConfig.Env)
	assert.Equal(t, "secrets", testConfig.Secrets)
}

func TestConfig_Template(t *testing.T) {
	t.Parallel()

	var testConfig struct {
		Path string
	}

	base := `
[group1]
path = "{{ .ProjectDir }}"
`

	dir := cconfigtest.SetupDirWithConfigs(t, base)

	config, err := cconfig.New(dir, projectDir, envTest)
	assert.NoError(t, err)

	err = config.Load("group1", &testConfig)
	assert.NoError(t, err)

	assert.Equal(t, projectDir, testConfig.Path)
}
