package cconfig_test

import (
	"io/ioutil"
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

	dir := cconfigtest.SetupDir(t)

	err := ioutil.WriteFile(path.Join(dir, "base.toml"), []byte(""), os.ModePerm)
	assert.NoError(t, err)

	err = ioutil.WriteFile(path.Join(dir, "test.toml"), []byte(""), os.ModePerm)
	assert.NoError(t, err)

	config, err := cconfig.NewConfig(dir, envTest)

	assert.NoError(t, err)
	assert.NotNil(t, config)
}

func TestNewConfig_MissingBase(t *testing.T) {
	t.Parallel()

	dir := cconfigtest.SetupDir(t)

	_, err := cconfig.NewConfig(dir, envTest)

	assert.Error(t, err)
}

func TestNewConfig_MissingEnv(t *testing.T) {
	t.Parallel()

	dir := cconfigtest.SetupDir(t)

	err := ioutil.WriteFile(path.Join(dir, "base.toml"), []byte(""), os.ModePerm)
	assert.NoError(t, err)

	_, err = cconfig.NewConfig(dir, envTest)

	assert.Error(t, err)
}

func TestConfig_Value(t *testing.T) {
	t.Parallel()

	dir := cconfigtest.SetupDir(t)

	err := ioutil.WriteFile(path.Join(dir, "base.toml"), []byte("file = \"base\""), os.ModePerm)
	assert.NoError(t, err)

	err = ioutil.WriteFile(path.Join(dir, "test.toml"), []byte(""), os.ModePerm)
	assert.NoError(t, err)

	config, err := cconfig.NewConfig(dir, envTest)
	assert.NoError(t, err)

	val, ok := config.Value("file").(string)

	assert.True(t, ok)
	assert.Equal(t, "base", val)
}

func TestConfig_Value_Override(t *testing.T) {
	t.Parallel()

	dir := cconfigtest.SetupDir(t)

	err := ioutil.WriteFile(path.Join(dir, "base.toml"), []byte("file = \"base\""), os.ModePerm)
	assert.NoError(t, err)

	err = ioutil.WriteFile(path.Join(dir, "test.toml"), []byte("file = \"test\""), os.ModePerm)
	assert.NoError(t, err)

	config, err := cconfig.NewConfig(dir, envTest)
	assert.NoError(t, err)

	val, ok := config.Value("file").(string)

	assert.True(t, ok)
	assert.Equal(t, "test", val)
}
