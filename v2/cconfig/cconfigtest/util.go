package cconfigtest

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// SetupDir creates an empty temp directory that can store config files.
// The directory is cleaned up after test run.
func SetupDir(t *testing.T) string {
	t.Helper()

	dir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, os.RemoveAll(dir))
	})

	return dir
}
