package console_test

import (
	"bytes"
	"errors"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tusharsoni/copper/v2/clogger"
	"github.com/tusharsoni/copper/v2/clogger/console"
)

func TestNew(t *testing.T) {
	t.Parallel()

	logger := console.New()

	_, ok := logger.(clogger.Logger)

	assert.NotNil(t, logger)
	assert.True(t, ok)
}

func TestLogger_Debug(t *testing.T) { //nolint:paralleltest
	var (
		buf    bytes.Buffer
		logger = console.New()
	)

	log.SetOutput(&buf)

	logger.Debug("test debug log")

	assert.Contains(t, buf.String(), "[DEBUG] test debug log")
}

func TestLogger_WithTags_Debug(t *testing.T) { //nolint:paralleltest
	var (
		buf    bytes.Buffer
		logger = console.New()
	)

	log.SetOutput(&buf)

	logger.
		WithTags(map[string]interface{}{
			"key": "val",
		}).
		WithTags(map[string]interface{}{
			"key2": "val2",
		}).Debug("test debug log")

	assert.Contains(t, buf.String(), "[DEBUG] test debug log where key=val,key2=val2")
}

func TestLogger_Info(t *testing.T) { //nolint:paralleltest
	var (
		buf    bytes.Buffer
		logger = console.New()
	)

	log.SetOutput(&buf)

	logger.Info("test info log")

	assert.Contains(t, buf.String(), "[INFO] test info log")
}

func TestLogger_Warn(t *testing.T) { //nolint:paralleltest
	var (
		buf    bytes.Buffer
		logger = console.New()
	)

	log.SetOutput(&buf)

	logger.Warn("test warn log", errors.New("test-error")) //nolint:goerr113

	assert.Contains(t, buf.String(), "[WARN] test warn log because\n> test-error")
}

func TestLogger_Error(t *testing.T) { //nolint:paralleltest
	var (
		buf    bytes.Buffer
		logger = console.New()
	)

	log.SetOutput(&buf)

	logger.Error("test error log", errors.New("test-error")) //nolint:goerr113

	assert.Contains(t, buf.String(), "[ERROR] test error log because\n> test-error")
}