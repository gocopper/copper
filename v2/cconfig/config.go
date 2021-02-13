package cconfig

import (
	"path"
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/tusharsoni/copper/v2/cerrors"
)

// Config provides methods to read app config.
type Config interface {
	Value(path string) interface{}
}

const (
	baseTomlConfigFileName = "base.toml"
	tomlExt                = ".toml"
)

// NewConfig provides an implementation of Config that reads config files in the
// dir. By default, it reads from base.toml and can be overridden by a file
// corresponding to the env. For 'test' env, the file should be test.toml.
func NewConfig(dir, env string) (Config, error) {
	baseConfigPath := path.Join(dir, baseTomlConfigFileName)
	envConfigPath := path.Join(dir, strings.ToLower(env)+tomlExt)

	baseTree, err := toml.LoadFile(baseConfigPath)
	if err != nil {
		return nil, cerrors.New(err, "failed to load base config file", map[string]interface{}{
			"path": baseConfigPath,
		})
	}

	envTree, err := toml.LoadFile(envConfigPath)
	if err != nil {
		return nil, cerrors.New(err, "failed to load env config file", map[string]interface{}{
			"env":  env,
			"path": envConfigPath,
		})
	}

	return &config{
		base: baseTree,
		env:  envTree,
	}, nil
}

type config struct {
	base *toml.Tree
	env  *toml.Tree
}

func (c *config) Value(path string) interface{} {
	val := c.env.Get(path)

	if val != nil {
		return val
	}

	return c.base.Get(path)
}
