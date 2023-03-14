package cconfigtest

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

// SetupDirWithConfigs creates a temp directory that can store config files.
// The directory is cleaned up after test run.
func SetupDirWithConfigs(t *testing.T, configs map[string]string) string {
	t.Helper()

	dir, err := os.MkdirTemp("", "")
	assert.NoError(t, err)

	for fp, data := range configs {
		err = os.WriteFile(path.Join(dir, fp), []byte(data), os.ModePerm)
		assert.NoError(t, err)
	}

	t.Cleanup(func() {
		assert.NoError(t, os.RemoveAll(dir))
	})

	return dir
}
