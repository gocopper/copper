package chttptest

import (
	"embed"
)

// HTMLDir embeds a directory that can be used with chttp.ReaderWriter
//go:embed html
var HTMLDir embed.FS
