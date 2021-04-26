package clogger_test

import (
	"testing"

	"github.com/gocopper/copper/clogger"
	"github.com/stretchr/testify/assert"
)

func TestLevel_String(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "DEBUG", clogger.LevelDebug.String())
	assert.Equal(t, "INFO", clogger.LevelInfo.String())
	assert.Equal(t, "WARN", clogger.LevelWarn.String())
	assert.Equal(t, "ERROR", clogger.LevelError.String())
	assert.Equal(t, "UNKNOWN", clogger.Level(-99).String())
}
