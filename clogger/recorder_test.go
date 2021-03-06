package clogger_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tusharsoni/copper/clogger"
)

func TestNewRecorder(t *testing.T) {
	t.Parallel()

	logger := clogger.NewRecorder(nil)

	_, ok := logger.(clogger.Logger)

	assert.NotNil(t, logger)
	assert.True(t, ok)
}

func TestRecorder_Debug(t *testing.T) {
	t.Parallel()

	var (
		logs   = make([]clogger.RecordedLog, 0)
		logger = clogger.NewRecorder(&logs)
	)

	logger.Debug("test debug log")

	assert.Len(t, logs, 1)

	log := logs[0]

	assert.Equal(t, clogger.LevelDebug, log.Level)
	assert.Equal(t, "test debug log", log.Msg)
	assert.Empty(t, log.Tags)
	assert.Nil(t, log.Error)
}

func TestRecorder_WithTags_Debug(t *testing.T) {
	t.Parallel()

	var (
		logs   = make([]clogger.RecordedLog, 0)
		logger = clogger.NewRecorder(&logs)
	)

	logger.
		WithTags(map[string]interface{}{
			"key": "val",
		}).
		WithTags(map[string]interface{}{
			"key2": "val2",
		}).Debug("test debug log")

	assert.Len(t, logs, 1)

	log := logs[0]

	assert.Equal(t, clogger.LevelDebug, log.Level)
	assert.Equal(t, "test debug log", log.Msg)
	assert.Equal(t, map[string]interface{}{
		"key":  "val",
		"key2": "val2",
	}, log.Tags)
	assert.Nil(t, log.Error)
}

func TestRecorder_Info(t *testing.T) {
	t.Parallel()

	var (
		logs   = make([]clogger.RecordedLog, 0)
		logger = clogger.NewRecorder(&logs)
	)

	logger.Info("test info log")

	assert.Len(t, logs, 1)

	log := logs[0]

	assert.Equal(t, clogger.LevelInfo, log.Level)
	assert.Equal(t, "test info log", log.Msg)
	assert.Empty(t, log.Tags)
	assert.Nil(t, log.Error)
}

func TestRecorder_Warn(t *testing.T) {
	t.Parallel()

	var (
		logs   = make([]clogger.RecordedLog, 0)
		logger = clogger.NewRecorder(&logs)
	)

	logger.Warn("test warn log", errors.New("test-err")) //nolint:goerr113

	assert.Len(t, logs, 1)

	log := logs[0]

	assert.Equal(t, clogger.LevelWarn, log.Level)
	assert.Equal(t, "test warn log", log.Msg)
	assert.Empty(t, log.Tags)
	assert.EqualError(t, log.Error, "test-err")
}

func TestRecorder_Error(t *testing.T) {
	t.Parallel()

	var (
		logs   = make([]clogger.RecordedLog, 0)
		logger = clogger.NewRecorder(&logs)
	)

	logger.Error("test error log", errors.New("test-err")) //nolint:goerr113

	assert.Len(t, logs, 1)

	log := logs[0]

	assert.Equal(t, clogger.LevelError, log.Level)
	assert.Equal(t, "test error log", log.Msg)
	assert.Empty(t, log.Tags)
	assert.EqualError(t, log.Error, "test-err")
}
