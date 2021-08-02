package chtmltest

import (
	"testing"

	"github.com/gocopper/copper/cconfig/cconfigtest"
	"github.com/gocopper/copper/chtml"
	"github.com/gocopper/copper/clogger"
	"github.com/stretchr/testify/assert"
)

// NewReaderWriter creates a *chtml.ReaderWriter suitable for use in tests
func NewReaderWriter(t *testing.T) *chtml.ReaderWriter {
	t.Helper()

	r, err := chtml.NewRenderer(chtml.NewRendererParams{
		HTMLDir:   HTMLDir,
		StaticDir: nil,
		AppConfig: cconfigtest.NewEmptyConfig(t),
	})
	assert.NoError(t, err)

	return chtml.NewReaderWriter(r, cconfigtest.NewEmptyConfig(t), clogger.NewNoop())
}
