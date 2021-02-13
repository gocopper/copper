package clogger_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tusharsoni/copper/v2/clogger"
)

func TestNewConsole(t *testing.T) {
	t.Parallel()

	logger := clogger.NewConsole()

	_, ok := logger.(clogger.Logger)

	assert.NotNil(t, logger)
	assert.True(t, ok)
}

func TestNewConsoleWithParams(t *testing.T) {
	t.Parallel()

	logger := clogger.NewConsoleWithParams(nil, nil)

	_, ok := logger.(clogger.Logger)

	assert.NotNil(t, logger)
	assert.True(t, ok)
}

func TestConsoleLogger_Debug(t *testing.T) {
	t.Parallel()

	var (
		buf    bytes.Buffer
		logger = clogger.NewConsoleWithParams(&buf, &buf)
	)

	logger.Debug("test debug log")

	assert.Contains(t, buf.String(), "[DEBUG] test debug log")
}

func TestConsoleLogger_WithTags_Debug(t *testing.T) {
	t.Parallel()

	var (
		buf    bytes.Buffer
		logger = clogger.NewConsoleWithParams(&buf, &buf)
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

func TestConsoleLogger_Info(t *testing.T) {
	t.Parallel()

	var (
		buf    bytes.Buffer
		logger = clogger.NewConsoleWithParams(&buf, &buf)
	)

	logger.Info("test info log")

	assert.Contains(t, buf.String(), "[INFO] test info log")
}

func TestConsoleLogger_Warn(t *testing.T) {
	t.Parallel()

	var (
		buf    bytes.Buffer
		logger = clogger.NewConsoleWithParams(&buf, &buf)
	)

	logger.Warn("test warn log", errors.New("test-error")) //nolint:goerr113

	assert.Contains(t, buf.String(), "[WARN] test warn log because\n> test-error")
}

func TestConsoleLogger_Error(t *testing.T) {
	t.Parallel()

	var (
		buf    bytes.Buffer
		logger = clogger.NewConsoleWithParams(&buf, &buf)
	)

	logger.Error("test error log", errors.New("test-error")) //nolint:goerr113

	assert.Contains(t, buf.String(), "[ERROR] test error log because\n> test-error")
}
