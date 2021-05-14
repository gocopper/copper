package chttp

import (
	"io/fs"
	"net/http"

	"github.com/gocopper/copper/clogger"

	"github.com/gocopper/copper/cconfig"
	"github.com/gocopper/copper/cerrors"
)

type (
	// StaticDir represents a directory that holds static resources (JS, CSS, images, etc.)
	StaticDir fs.FS

	// HTMLRouter provides a route that exposes the StaticDir and handles server-side HTML components
	HTMLRouter struct {
		dir    fs.FS
		rw     *ReaderWriter
		html   *HTMLRenderer
		logger clogger.Logger
	}

	// NewHTMLRouterParams holds the params needed to instantiate a new HTMLRouter
	NewHTMLRouterParams struct {
		StaticDir StaticDir
		RW        *ReaderWriter
		HTML      *HTMLRenderer
		AppConfig cconfig.Config
		Logger    clogger.Logger
	}
)

// NewHTMLRouter instantiates a new HTMLRouter
func NewHTMLRouter(p NewHTMLRouterParams) (*HTMLRouter, error) {
	var config config

	router := HTMLRouter{
		dir:    p.StaticDir,
		rw:     p.RW,
		html:   p.HTML,
		logger: p.Logger,
	}

	err := p.AppConfig.Load("chttp", &config)
	if err != nil {
		return nil, cerrors.New(err, "failed to load chttp config", nil)
	}

	return &router, nil
}

// Routes defines the /static/ route that serves the StaticDir
func (ro *HTMLRouter) Routes() []Route {
	return []Route{
		{
			Path:    "/static/{path:.*}",
			Methods: []string{http.MethodGet},
			Handler: http.FileServer(http.FS(ro.dir)).ServeHTTP,
		},
	}
}
