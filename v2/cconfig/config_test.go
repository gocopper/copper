package cconfig_test

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tusharsoni/copper/v2/cconfig"
	"github.com/tusharsoni/copper/v2/cconfig/cconfigtest"
)

const (
	envTest = "test"
)

func TestNewConfig(t *testing.T) {
	t.Parallel()

	configDir := cconfigtest.SetupDirWithConfigs(t, "", "")

	config, err := cconfig.New(configDir, envTest)

	assert.NoError(t, err)
	assert.NotNil(t, config)
}

func TestNewConfig_MissingBase(t *testing.T) {
	t.Parallel()

	dir := cconfigtest.SetupDirWithConfigs(t, "", "")

	assert.NoError(t, os.Remove(path.Join(string(dir), "base.toml")))

	_, err := cconfig.New(dir, envTest)

	assert.Error(t, err)
}

func TestNewConfig_MissingEnv(t *testing.T) {
	t.Parallel()

	dir := cconfigtest.SetupDirWithConfigs(t, "", "")

	assert.NoError(t, os.Remove(path.Join(string(dir), "test.toml")))

	_, err := cconfig.New(dir, envTest)

	assert.Error(t, err)
}

func TestConfig_Load_Default(t *testing.T) {
	t.Parallel()

	var testConfig struct {
		Value string `default:"default"`
	}

	dir := cconfigtest.SetupDirWithConfigs(t, "", "")

	config, err := cconfig.New(dir, envTest)
	assert.NoError(t, err)

	err = config.Load("group1", &testConfig)
	assert.NoError(t, err)

	assert.Equal(t, "default", testConfig.Value)
}

func TestConfig_Load_Base(t *testing.T) {
	t.Parallel()

	var testConfig struct {
		Value string `default:"default"`
	}

	base := `
[group1]
value = "base"
`

	dir := cconfigtest.SetupDirWithConfigs(t, base, "")

	config, err := cconfig.New(dir, envTest)
	assert.NoError(t, err)

	err = config.Load("group1", &testConfig)
	assert.NoError(t, err)

	assert.Equal(t, "base", testConfig.Value)
}

func TestConfig_Load_Env(t *testing.T) {
	t.Parallel()

	var testConfig struct {
		Value string `default:"default"`
	}

	base := `
[group1]
value = "base"
`

	env := `
[group1]
value = "env"`

	dir := cconfigtest.SetupDirWithConfigs(t, base, env)

	config, err := cconfig.New(dir, envTest)
	assert.NoError(t, err)

	err = config.Load("group1", &testConfig)
	assert.NoError(t, err)

	assert.Equal(t, "env", testConfig.Value)
}
