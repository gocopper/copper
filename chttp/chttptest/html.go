package chttptest

import (
	"embed"
)

// HTMLDir embeds a directory that can be used with chttp.ReaderWriter
//go:embed src
var HTMLDir embed.FS
