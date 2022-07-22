package clogger_test

import (
	"testing"

	"github.com/gocopper/copper/clogger"
	"github.com/stretchr/testify/assert"
)

func TestNewNoop(t *testing.T) {
	t.Parallel()

	logger := clogger.NewNoop()
	assert.NotNil(t, logger)
}

func TestNoopLogger_WithTags(t *testing.T) {
	t.Parallel()

	logger := clogger.NewNoop().WithTags(nil)

	assert.NotNil(t, logger)
}

func TestNoopLogger_Debug(t *testing.T) {
	t.Parallel()

	logger := clogger.NewNoop().WithTags(nil)

	logger.Debug("test-debug")
}

func TestNoopLogger_Info(t *testing.T) {
	t.Parallel()

	logger := clogger.NewNoop().WithTags(nil)

	logger.Info("info")
}

func TestNoopLogger_Warn(t *testing.T) {
	t.Parallel()

	logger := clogger.NewNoop().WithTags(nil)

	logger.Warn("warn", nil)
}

func TestNoopLogger_Error(t *testing.T) {
	t.Parallel()

	logger := clogger.NewNoop().WithTags(nil)

	logger.Error("error", nil)
}
