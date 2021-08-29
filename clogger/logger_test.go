package clogger_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/gocopper/copper/clogger"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	logger := clogger.New()

	_, ok := logger.(clogger.Logger)

	assert.NotNil(t, logger)
	assert.True(t, ok)
}

func TestNewWithConfig(t *testing.T) {
	t.Parallel()

	log, err := ioutil.TempFile("", "*")
	assert.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, os.Remove(log.Name()))
	})

	logger, err := clogger.NewWithConfig(clogger.Config{
		Out:    log.Name(),
		Err:    log.Name(),
		Format: clogger.FormatPlain,
	})
	assert.NoError(t, err)

	_, ok := logger.(clogger.Logger)

	assert.NotNil(t, logger)
	assert.True(t, ok)
}

func TestNewWithConfig_OutFileErr(t *testing.T) {
	t.Parallel()

	log, err := ioutil.TempFile("", "*")
	assert.NoError(t, err)

	assert.NoError(t, os.Chmod(log.Name(), 0000))

	t.Cleanup(func() {
		assert.NoError(t, os.Remove(log.Name()))
	})

	_, err = clogger.NewWithConfig(clogger.Config{
		Out:    log.Name(),
		Format: clogger.FormatPlain,
	})
	assert.Error(t, err)
}

func TestNewWithConfig_ErrFileErr(t *testing.T) {
	t.Parallel()

	log, err := ioutil.TempFile("", "*")
	assert.NoError(t, err)

	assert.NoError(t, os.Chmod(log.Name(), 0000))

	t.Cleanup(func() {
		assert.NoError(t, os.Remove(log.Name()))
	})

	_, err = clogger.NewWithConfig(clogger.Config{
		Err:    log.Name(),
		Format: clogger.FormatPlain,
	})
	assert.Error(t, err)
}

func TestNewWithParams(t *testing.T) {
	t.Parallel()

	logger := clogger.NewWithWriters(nil, nil, clogger.FormatPlain)

	_, ok := logger.(clogger.Logger)

	assert.NotNil(t, logger)
	assert.True(t, ok)
}

func TestLogger_Debug(t *testing.T) {
	t.Parallel()

	var (
		buf    bytes.Buffer
		logger = clogger.NewWithWriters(&buf, &buf, clogger.FormatPlain)
	)

	logger.Debug("test debug log")

	assert.Contains(t, buf.String(), "[DEBUG] test debug log")
}

func TestLogger_WithTags_Debug(t *testing.T) {
	t.Parallel()

	var (
		buf    bytes.Buffer
		logger = clogger.NewWithWriters(&buf, &buf, clogger.FormatPlain)
	)

	logger.
		WithTags(map[string]interface{}{
			"key": "val",
		}).
		WithTags(map[string]interface{}{
			"key2": "val2",
		}).Debug("test debug log")

	assert.Contains(t, buf.String(), "[DEBUG] test debug log where key=val,key2=val2")
}

func TestLogger_Info(t *testing.T) {
	t.Parallel()

	var (
		buf    bytes.Buffer
		logger = clogger.NewWithWriters(&buf, &buf, clogger.FormatPlain)
	)

	logger.Info("test info log")

	assert.Contains(t, buf.String(), "[INFO] test info log")
}

func TestLogger_Warn(t *testing.T) {
	t.Parallel()

	var (
		buf    bytes.Buffer
		logger = clogger.NewWithWriters(&buf, &buf, clogger.FormatPlain)
	)

	logger.Warn("test warn log", errors.New("test-error")) //nolint:goerr113

	assert.Contains(t, buf.String(), "[WARN] test warn log because\n> test-error")
}

func TestLogger_Error(t *testing.T) {
	t.Parallel()

	var (
		buf    bytes.Buffer
		logger = clogger.NewWithWriters(&buf, &buf, clogger.FormatPlain)
	)

	logger.Error("test error log", errors.New("test-error")) //nolint:goerr113

	assert.Contains(t, buf.String(), "[ERROR] test error log because\n> test-error")
}
