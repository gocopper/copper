package chttptest

import (
	"embed"
	"testing"

	"github.com/gocopper/copper/chttp"
	"github.com/gocopper/copper/clogger"
	"github.com/stretchr/testify/assert"
)

// HTMLDir embeds a directory that can be used with chttp.ReaderWriter
//go:embed src
var HTMLDir embed.FS

// NewReaderWriter creates a *chttp.ReaderWriter suitable for use in tests
func NewReaderWriter(t *testing.T) *chttp.ReaderWriter {
	t.Helper()

	r, err := chttp.NewHTMLRenderer(chttp.NewHTMLRendererParams{
		HTMLDir:   HTMLDir,
		StaticDir: nil,
		Config:    chttp.Config{},
		Logger:    clogger.NewNoop(),
	})
	assert.NoError(t, err)

	return chttp.NewReaderWriter(r, chttp.Config{}, clogger.NewNoop())
}
