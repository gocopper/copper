package noop_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tusharsoni/copper/v2/clogger"
	"github.com/tusharsoni/copper/v2/clogger/noop"
)

func TestNew(t *testing.T) {
	t.Parallel()

	logger := noop.New()

	_, ok := logger.(clogger.Logger)

	assert.NotNil(t, logger)
	assert.True(t, ok)
}

func TestLogger_WithTags(t *testing.T) {
	t.Parallel()

	logger := noop.New().WithTags(nil)

	assert.NotNil(t, logger)
}
