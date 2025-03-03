package chttptest

import (
	"testing"

	"github.com/gocopper/copper/chttp"
	"github.com/gocopper/copper/clogger"
)

// NewJSONReaderWriter creates a *chttp.NewJSONReaderWriter suitable for use in tests
func NewJSONReaderWriter(t *testing.T) *chttp.JSONReaderWriter {
	t.Helper()

	return chttp.NewJSONReaderWriter(chttp.Config{}, clogger.NewNoop())
}
