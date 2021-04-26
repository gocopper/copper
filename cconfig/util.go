package cconfig

import (
	"bytes"
	"path/filepath"
	"text/template"

	"github.com/gocopper/copper/cerrors"
	"github.com/pelletier/go-toml"
)

func loadTOMLTemplate(fp string, pd ProjectDir) (*toml.Tree, error) {
	var buf bytes.Buffer

	apd, err := filepath.Abs(string(pd))
	if err != nil {
		return nil, cerrors.New(err, "failed to get absolute project dir path", map[string]interface{}{
			"projectDir": pd,
		})
	}

	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		return nil, cerrors.New(err, "failed to parse template file", map[string]interface{}{
			"path": fp,
		})
	}

	err = tmpl.Execute(&buf, map[string]interface{}{
		"ProjectDir": apd,
	})
	if err != nil {
		return nil, cerrors.New(err, "failed to execute config template", map[string]interface{}{
			"path":       fp,
			"projectDir": apd,
		})
	}

	tree, err := toml.LoadBytes(buf.Bytes())
	if err != nil {
		return nil, cerrors.New(err, "failed to load toml from template bytes", map[string]interface{}{
			"path": fp,
		})
	}

	return tree, nil
}
