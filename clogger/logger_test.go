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

	logger, err := clogger.New(clogger.Config{}, nil)
	assert.NoError(t, err)
	assert.NotNil(t, logger)
}

func TestNewCore(t *testing.T) {
	t.Parallel()

	logger, err := clogger.NewCore(clogger.Config{})
	assert.NoError(t, err)
	assert.NotNil(t, logger)
}

func TestNew_WithConfig(t *testing.T) {
	t.Parallel()

	log, err := os.CreateTemp("", "*")
	assert.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, os.Remove(log.Name()))
	})

	logger, err := clogger.New(clogger.Config{
		Out:    log.Name(),
		Err:    log.Name(),
		Format: clogger.FormatPlain,
	}, nil)
	assert.NoError(t, err)
	assert.NotNil(t, logger)
}

func TestNew_OutFileErr(t *testing.T) {
	t.Parallel()

	log, err := os.CreateTemp("", "*")
	assert.NoError(t, err)

	assert.NoError(t, os.Chmod(log.Name(), 0000))

	t.Cleanup(func() {
		assert.NoError(t, os.Remove(log.Name()))
	})

	_, err = clogger.New(clogger.Config{
		Out:    log.Name(),
		Format: clogger.FormatPlain,
	}, nil)
	assert.Error(t, err)
}

func TestNew_ErrFileErr(t *testing.T) {
	t.Parallel()

	log, err := os.CreateTemp("", "*")
	assert.NoError(t, err)

	assert.NoError(t, os.Chmod(log.Name(), 0000))

	t.Cleanup(func() {
		assert.NoError(t, os.Remove(log.Name()))
	})

	_, err = clogger.New(clogger.Config{
		Err:    log.Name(),
		Format: clogger.FormatPlain,
	}, nil)
	assert.Error(t, err)
}

func TestLogger_Debug(t *testing.T) {
	t.Parallel()

	var (
		buf    bytes.Buffer
		logger = clogger.NewWithWriters(&buf, &buf, clogger.FormatPlain, nil, nil, nil)
	)

	logger.Debug("test debug log")

	assert.Contains(t, buf.String(), "[DEBUG] test debug log")
}

func TestLogger_WithTags_Debug(t *testing.T) {
	t.Parallel()

	var (
		buf    bytes.Buffer
		logger = clogger.NewWithWriters(&buf, &buf, clogger.FormatPlain, nil, nil, nil)
	)

	logger.
		WithTags(map[string]any{
			"key": "val",
		}).
		WithTags(map[string]any{
			"key2": "val2",
		}).Debug("test debug log")

	assert.Contains(t, buf.String(), "[DEBUG] test debug log where key2=val2,key=val")
}

func TestLogger_WithTags_RedactedFields(t *testing.T) {
	t.Parallel()

	for _, format := range []clogger.Format{clogger.FormatJSON, clogger.FormatPlain} {
		var (
			buf    bytes.Buffer
			logger = clogger.NewWithWriters(&buf, &buf, format, []string{
				"secret", "password", "userPin",
			}, nil, nil)

			testErr = cerrors.New(nil, "test-error", map[string]any{
				"secret":   "my_api_key",
				"user-pin": "12456",
				"data": map[string]string{
					"password": "abc123",
				},
			})
		)

		logger.WithTags(map[string]any{
			"passwordOwner": "abc123",
			"USER_PIN":      "123456",
			"params": map[string]string{
				"myPassword": "abc123",
			},
		}).Error("test debug log", testErr)

		assert.NotContains(t, buf.String(), "my_api_key")
		assert.NotContains(t, buf.String(), "12456")
		assert.NotContains(t, buf.String(), "123456")
		assert.NotContains(t, buf.String(), "abc123")
		assert.Contains(t, buf.String(), "redact")
	}
}

func TestLogger_Info(t *testing.T) {
	t.Parallel()

	var (
		buf    bytes.Buffer
		logger = clogger.NewWithWriters(&buf, &buf, clogger.FormatPlain, nil, nil, nil)
	)

	logger.Info("test info log")

	assert.Contains(t, buf.String(), "[INFO] test info log", nil)
}

func TestLogger_Warn(t *testing.T) {
	t.Parallel()

	var (
		buf    bytes.Buffer
		logger = clogger.NewWithWriters(&buf, &buf, clogger.FormatPlain, nil, nil, nil)
	)

	logger.Warn("test warn log", errors.New("test-error")) //nolint:goerr113

	assert.Contains(t, buf.String(), "[WARN] test warn log because\n> test-error")
}

func TestLogger_Error(t *testing.T) {
	t.Parallel()

	var (
		buf    bytes.Buffer
		logger = clogger.NewWithWriters(&buf, &buf, clogger.FormatPlain, nil, nil, nil)
	)

	logger.Error("test error log", errors.New("test-error")) //nolint:goerr113

	assert.Contains(t, buf.String(), "[ERROR] test error log because\n> test-error")
}
