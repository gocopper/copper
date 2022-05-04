package cconfigtest

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

// SetupDirWithConfigs creates a temp directory that can store config files.
// The directory is cleaned up after test run.
func SetupDirWithConfigs(t *testing.T, configs map[string]string) string {
	t.Helper()

	dir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)

	for fp, data := range configs {
		err = ioutil.WriteFile(path.Join(dir, fp), []byte(data), os.ModePerm)
		assert.NoError(t, err)
	}

	t.Cleanup(func() {
		assert.NoError(t, os.RemoveAll(dir))
	})

	return dir
}
