package cconfig

import (
	"path"
	"reflect"
	"strings"

	"github.com/imdario/mergo"
	"github.com/pelletier/go-toml"
	"github.com/tusharsoni/copper/cerrors"
)

// Env defines the various environments the app can be configured for.
// The Env may be dev, test, staging, or prod.
type Env string

// Dir defines the directory where config file(s) live.
type Dir string

// Config provides methods to read app config.
type Config interface {
	Load(key string, dest interface{}) error
}

const (
	baseTomlConfigFileName = "base.toml"
	tomlExt                = ".toml"
)

// New provides an implementation of Config that reads config files in the
// dir. By default, it reads from base.toml and can be overridden by a file
// corresponding to the env. For 'test' env, the file should be test.toml.
func New(dir Dir, env Env) (Config, error) {
	baseConfigPath := path.Join(string(dir), baseTomlConfigFileName)
	envConfigPath := path.Join(string(dir), strings.ToLower(string(env))+tomlExt)

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

func (c *config) Load(key string, dest interface{}) error {
	var (
		base = &toml.Tree{}
		env  = &toml.Tree{}
	)

	if c.base.Has(key) {
		base = c.base.Get(key).(*toml.Tree)
	}

	if c.env.Has(key) {
		env = c.env.Get(key).(*toml.Tree)
	}

	// create a new value with the same type as dest
	// we will unmarshal with empty config to set all the
	// default values as set on the dest's struct tags
	defaults := reflect.New(reflect.TypeOf(dest).Elem()).Interface()

	err := toml.Unmarshal([]byte(""), defaults)
	if err != nil {
		return cerrors.New(err, "failed to load config defaults", map[string]interface{}{
			"key": key,
		})
	}

	err = env.Unmarshal(dest)
	if err != nil {
		return cerrors.New(err, "failed to unmarshal env config", map[string]interface{}{
			"key": key,
		})
	}

	// removes default values from dest by 'merging' dest + defaults
	// using a custom transformer. the transformer checks if the dest
	// has a default value. if so, it sets it to its zero value.
	err = mergo.Merge(dest, defaults, mergo.WithTransformers(c))
	if err != nil {
		return cerrors.New(err, "failed to remove default values from config", map[string]interface{}{
			"key": key,
		})
	}

	baseDest := reflect.New(reflect.TypeOf(dest).Elem()).Interface()

	err = base.Unmarshal(baseDest)
	if err != nil {
		return cerrors.New(err, "failed to unmarshal base config", map[string]interface{}{
			"key": key,
		})
	}

	err = mergo.Merge(dest, baseDest)
	if err != nil {
		return cerrors.New(err, "failed to merge env with base config", map[string]interface{}{
			"key": key,
		})
	}

	return nil
}

func (c *config) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	return func(dst, src reflect.Value) error {
		if dst.Interface() == src.Interface() {
			dst.Set(reflect.Zero(dst.Type()))
		}

		return nil
	}
}
