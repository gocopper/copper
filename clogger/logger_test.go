package clogger_test

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/gocopper/copper/cerrors"

	"github.com/gocopper/copper/clogger"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	logger := clogger.New()
	assert.NotNil(t, logger)
}

func TestNewWithConfig(t *testing.T) {
	t.Parallel()

	log, err := os.CreateTemp("", "*")
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
	assert.NotNil(t, logger)
}

func TestNewWithConfig_OutFileErr(t *testing.T) {
	t.Parallel()

	log, err := os.CreateTemp("", "*")
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

	log, err := os.CreateTemp("", "*")
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

	logger := clogger.NewWithWriters(nil, nil, clogger.FormatPlain, nil)
	assert.NotNil(t, logger)
}

func TestLogger_Debug(t *testing.T) {
	t.Parallel()

	var (
		buf    bytes.Buffer
		logger = clogger.NewWithWriters(&buf, &buf, clogger.FormatPlain, nil)
	)

	logger.Debug("test debug log")

	assert.Contains(t, buf.String(), "[DEBUG] test debug log")
}

func TestLogger_WithTags_Debug(t *testing.T) {
	t.Parallel()

	var (
		buf    bytes.Buffer
		logger = clogger.NewWithWriters(&buf, &buf, clogger.FormatPlain, nil)
	)

	logger.
		WithTags(map[string]interface{}{
			"key": "val",
		}).
		WithTags(map[string]interface{}{
			"key2": "val2",
		}).Debug("test debug log")

	assert.Contains(t, buf.String(), "[DEBUG] test debug log where key2=val2,key=val")
}

func TestLogger_WithTags_RedactedFields(t *testing.T) {
	t.Parallel()

	var (
		buf    bytes.Buffer
		logger = clogger.NewWithWriters(&buf, &buf, clogger.FormatPlain, []string{
			"secret", "password", "userPin",
		})

		testErr = cerrors.New(nil, "test-error", map[string]interface{}{
			"secret":   "my_api_key",
			"user-pin": "12456",
		})
	)

	logger.WithTags(map[string]interface{}{
		"password": "abc123",
		"USER_PIN": "123456",
	}).Error("test debug log", testErr)

	assert.NotContains(t, buf.String(), "my_api_key")
	assert.NotContains(t, buf.String(), "12456")
	assert.NotContains(t, buf.String(), "abc123")
	assert.Contains(t, buf.String(), "<redacted>")
}

func TestLogger_Info(t *testing.T) {
	t.Parallel()

	var (
		buf    bytes.Buffer
		logger = clogger.NewWithWriters(&buf, &buf, clogger.FormatPlain, nil)
	)

	logger.Info("test info log")

	assert.Contains(t, buf.String(), "[INFO] test info log", nil)
}

func TestLogger_Warn(t *testing.T) {
	t.Parallel()

	var (
		buf    bytes.Buffer
		logger = clogger.NewWithWriters(&buf, &buf, clogger.FormatPlain, nil)
	)

	logger.Warn("test warn log", errors.New("test-error")) //nolint:goerr113

	assert.Contains(t, buf.String(), "[WARN] test warn log because\n> test-error")
}

func TestLogger_Error(t *testing.T) {
	t.Parallel()

	var (
		buf    bytes.Buffer
		logger = clogger.NewWithWriters(&buf, &buf, clogger.FormatPlain, nil)
	)

	logger.Error("test error log", errors.New("test-error")) //nolint:goerr113

	assert.Contains(t, buf.String(), "[ERROR] test error log because\n> test-error")
}
