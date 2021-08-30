package cconfig_test

import (
	"testing"

	"github.com/gocopper/copper/cconfig"
	"github.com/gocopper/copper/cconfig/cconfigtest"
	"github.com/stretchr/testify/assert"
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
		Local   string `default:"default"`
	}

	base := `
[group1]
base = "base"
`

	env := `
[group1]
env = "env"
`

	local := `
[group1]
local = "local"
`

	dir := cconfigtest.SetupDirWithConfigs(t, base, env, local)

	config, err := cconfig.New(dir, projectDir, envTest)
	assert.NoError(t, err)

	err = config.Load("group1", &testConfig)
	assert.NoError(t, err)

	assert.Equal(t, "default", testConfig.Default)
	assert.Equal(t, "base", testConfig.Base)
	assert.Equal(t, "env", testConfig.Env)
	assert.Equal(t, "local", testConfig.Local)
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
