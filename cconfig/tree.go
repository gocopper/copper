package cconfig

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/gocopper/copper/cerrors"
	"github.com/pelletier/go-toml"
)

//nolint:funlen
func loadTree(fp, overrides string, disableKeyOverrides bool) (*toml.Tree, error) {
	funcMap := template.FuncMap{
		"exec": execCmd,
	}

	tmpl, err := template.New(filepath.Base(fp)).Funcs(funcMap).ParseFiles(fp)
	if err != nil {
		return nil, cerrors.New(err, "failed to parse config file as template", map[string]interface{}{
			"path": fp,
		})
	}

	envVars := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		envVars[pair[0]] = pair[1]
	}

	var tomlOut strings.Builder
	err = tmpl.Execute(&tomlOut, map[string]interface{}{
		"EnvVars": envVars,
	})
	if err != nil {
		return nil, cerrors.New(err, "failed to execute config file template", map[string]interface{}{
			"path": fp,
		})
	}

	tree, err := toml.LoadBytes([]byte(tomlOut.String()))
	if err != nil {
		return nil, cerrors.New(err, "failed to load config file", map[string]interface{}{
			"path": fp,
		})
	}

	// If the TOML tree does not have a top-level 'extends' key, we can return the tree as-is
	if !tree.Has("extends") {
		return tree, nil
	}

	parentFilePaths := make([]string, 0)

	// The extends key can be a string or a list of strings representing the config file paths that need to be loaded
	switch extends := tree.Get("extends").(type) {
	case string:
		parentFilePaths = append(parentFilePaths, extends)

	// If extends is set to a list, verify each value is a valid string, and add it to parentFilePaths
	case []interface{}:
		for i := range extends {
			parentFilePath, ok := extends[i].(string)
			if !ok {
				return nil, cerrors.New(nil, "extends can only contain strings", map[string]interface{}{
					"path":  fp,
					"value": extends[i],
				})
			}

			parentFilePaths = append(parentFilePaths, parentFilePath)
		}
	default:
		return nil, cerrors.New(nil, "'extends' must be string or []string", map[string]interface{}{
			"path": fp,
			"type": reflect.TypeOf(extends).String(),
		})
	}

	// Load each parentFilePath in-order
	for _, parentFP := range parentFilePaths {
		parentFilePath := filepath.Join(filepath.Dir(fp), parentFP)

		// Load the parent tree at the given path defined by the extends key. Note that this is a recursive call
		// that will load all ancestors.
		parentTree, err := loadTree(parentFilePath, "", disableKeyOverrides)
		if err != nil {
			return nil, cerrors.New(err, "failed to load parent tree", map[string]interface{}{
				"parentPath": parentFilePath,
			})
		}

		// Once the parent tree and its ancestors are loaded, we need to merge it with our current tree
		tree, err = mergeTrees(parentTree, tree, disableKeyOverrides)
		if err != nil {
			return nil, cerrors.New(err, "failed to merge with parent tree", map[string]interface{}{
				"parentPath": parentFilePath,
			})
		}
	}

	// Apply overrides
	for _, ov := range strings.Split(overrides, ";") {
		t, err := toml.Load(ov)
		if err != nil {
			return nil, cerrors.New(err, "failed to parse override as TOML", map[string]interface{}{
				"override": ov,
			})
		}

		tree, err = mergeTrees(tree, t, disableKeyOverrides)
		if err != nil {
			return nil, cerrors.New(err, "failed to merge tree with overrides", map[string]interface{}{
				"override": ov,
			})
		}
	}
	return tree, nil
}

//nolint:funlen
func mergeTrees(base, override *toml.Tree, disableKeyOverrides bool) (*toml.Tree, error) {
	// For each key in the override tree, we need to apply it to the base tree
	for _, key := range override.Keys() {
		switch keyVal := override.Get(key).(type) {
		// If the value at the given key is a TOML tree (aka a table according to the spec), we need to merge it with
		// the base table.
		// For example, if the base tree contains:
		// [group1]
		// key1="val1"
		// and the override tree contains:
		// [group1]
		// key2="val"2
		// We need to load it in such a way where group1 contains both key1 and key2
		case *toml.Tree:
			// If the base does not contain the key, we can set the entire value from the override table as-is
			if !base.Has(key) {
				base.Set(key, keyVal)
				continue
			}

			// Verify that the value type in the base and override trees are the same. For example, this is invalid:
			// # base.toml
			// group1 = "I am a string"
			//
			// # prod.toml
			// extends = "base.toml"
			// [group1] # I am a table!
			// key1="val"1
			//
			// The above configuration is invalid because group1 is a table in prod.toml but a string in base.toml. As
			// a result, they cannot be merged.
			baseTree, ok := base.Get(key).(*toml.Tree)
			if !ok {
				return nil, cerrors.New(nil, "base and override key types don't match", map[string]interface{}{
					"key": key,
				})
			}

			// Now that we have two trees, we can merge them recursively
			mergedTree, err := mergeTrees(baseTree, keyVal, disableKeyOverrides)
			if err != nil {
				return nil, cerrors.New(err, "failed to merge tree for key", map[string]interface{}{
					"key": key,
				})
			}

			base.Set(key, mergedTree)

		// This handles all non-table keys. The key, as found in the override tree, is set on the base tree. If the base
		// tree already has a value for the key, it is only overridden if disableKeyOverrides is false. If a key is
		// being overridden with disableKeyOverrides=true, an error is returned.
		default:
			if base.Has(key) && disableKeyOverrides {
				return nil, cerrors.New(nil, "key is being overridden when key overrides are disabled", map[string]interface{}{
					"key": key,
				})
			}

			base.Set(key, keyVal)
		}
	}

	return base, nil
}

func execCmd(cmd string) (string, error) {
	c := exec.Command("sh", "-c", cmd)

	var stdout, stderr strings.Builder
	c.Stdout = &stdout
	c.Stderr = &stderr

	err := c.Run()
	if err != nil {
		return "", cerrors.New(err, "failed to execute command in config template", map[string]interface{}{
			"cmd":    cmd,
			"stdout": stdout.String(),
			"stderr": stderr.String(),
		})
	}

	return strings.TrimSpace(stdout.String()), nil
}
