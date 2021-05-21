package chttptest

import (
	"testing"

	"github.com/gocopper/copper/cconfig/cconfigtest"

	"github.com/gocopper/copper/chttp"
	"github.com/gocopper/copper/clogger"
	"github.com/stretchr/testify/assert"
)

// NewReaderWriter creates a *chttp.ReaderWriter suitable for use in tests
func NewReaderWriter(t *testing.T) *chttp.ReaderWriter {
	t.Helper()

	r, err := chttp.NewHTMLRenderer(chttp.NewHTMLRendererParams{
		HTMLDir:   HTMLDir,
		StaticDir: nil,
		AppConfig: cconfigtest.NewEmptyConfig(t),
	})
	assert.NoError(t, err)

	return chttp.NewReaderWriter(r, cconfigtest.NewEmptyConfig(t), clogger.NewNoop())
}
