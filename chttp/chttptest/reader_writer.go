package chttptest

import (
	"testing"

	"github.com/gocopper/copper/chttp"
	"github.com/gocopper/copper/clogger"
)

// NewReaderWriter creates a *chttp.ReaderWriter suitable for use in tests
func NewReaderWriter(t *testing.T) *chttp.ReaderWriter {
	t.Helper()

	return chttp.NewReaderWriter(clogger.NewNoop())
}
