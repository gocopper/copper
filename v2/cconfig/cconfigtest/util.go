package cconfigtest

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tusharsoni/copper/v2/cconfig"
)

// SetupDirWithConfigs creates a temp directory that can store config files.
// It creates base.toml and test.toml with the given strings.
// The directory is cleaned up after test run.
func SetupDirWithConfigs(t *testing.T, base, test string) cconfig.Dir {
	t.Helper()

	dir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)

	err = ioutil.WriteFile(path.Join(dir, "base.toml"), []byte(base), os.ModePerm)
	assert.NoError(t, err)

	err = ioutil.WriteFile(path.Join(dir, "test.toml"), []byte(test), os.ModePerm)
	assert.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, os.RemoveAll(dir))
	})

	return cconfig.Dir(dir)
}
