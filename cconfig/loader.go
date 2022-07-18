package cconfig

import (
	"github.com/gocopper/copper/cerrors"
	"github.com/pelletier/go-toml"
)

type (
	// Path defines the path to the config file.
	Path string
	// Overrides defines a ';' separated string of config overrides in TOML format
	Overrides string
)

// Loader provides methods to load config files into structs.
type Loader interface {
	// Load reads the config values defined under the given key (TOML table) and sets them into the dest struct.
	// For example:
	// # prod.toml
	// key1 = "val1"
	//
	// # config.go
	// type MyConfig struct {
	//   Key1 string `toml:"key1"`
	// }
	//
	// func LoadMyConfig(loader cconfig.Loader) (MyConfig, error) {
	//   var config MyConfig
	//
	//   err := loader.Load("my_config", &config)
	//   if err != nil {
	//     return MyConfig{}, cerrors.New(err, "failed to load configs for my_config", nil)
	//   }
	//
	//   return config, nil
	// }
	Load(key string, dest interface{}) error
}

// New provides an implementation of Loader that reads a config file at the given file path. It supports extending the
// config file at the given path by using an 'extends' key. For example, the config file may be extended like so:
//
// # base.toml
// key1 = "val1"
//
// # prod.toml
// extends = "base.toml
// key2 = "val2"
//
// If New is called with the path to prod.toml, it loads both key2 (from prod.toml) and key1 (from base.toml) since
// prod.toml extends base.toml.
//
// The extends key can support multiple files like so:
// extends = ["base.toml", "secrets.toml"]
//
// If a config key is present in multiple files, New returns an error. For example, if prod.toml sets a value for 'key1'
// that has already been set in base.toml, an error will be returned. To enable key overrides see NewWithKeyOverrides.
func New(fp Path, ov Overrides) (Loader, error) {
	return newLoader(string(fp), string(ov), true)
}

// NewWithKeyOverrides works exactly the same way as New except it supports key overrides. For example, this is a valid
// config:
// # base.toml
// key1 = "val1"
//
// # prod.toml
// extends = "base.toml
// key1 = "val2"
//
// If prod.toml is loaded, key1 will be set to "val2" since it has been overridden in prod.toml.
func NewWithKeyOverrides(fp Path, overrides Overrides) (Loader, error) {
	return newLoader(string(fp), string(overrides), false)
}

func newLoader(fp, overrides string, disableKeyOverrides bool) (*loader, error) {
	tree, err := loadTree(fp, overrides, disableKeyOverrides)
	if err != nil {
		return nil, cerrors.New(err, "failed to load config tree", map[string]interface{}{
			"path": fp,
		})
	}

	return &loader{
		tree: tree,
	}, nil
}

type loader struct {
	tree *toml.Tree
}

func (l *loader) Load(key string, dest interface{}) error {
	if !l.tree.Has(key) {
		return nil
	}

	keyTree, ok := l.tree.Get(key).(*toml.Tree)
	if !ok {
		return cerrors.New(nil, "invalid key type", map[string]interface{}{
			"key": key,
		})
	}

	err := keyTree.Unmarshal(dest)
	if err != nil {
		return cerrors.New(err, "failed to unmarshal config into dest", map[string]interface{}{
			"key": key,
		})
	}

	return nil
}
