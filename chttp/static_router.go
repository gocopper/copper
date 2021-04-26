package chttp

import (
	"io/fs"
	"net/http"
	"os"

	"github.com/gocopper/copper/cconfig"
	"github.com/gocopper/copper/cerrors"
)

type (
	// StaticDir represents a directory that holds static resources (JS, CSS, images, etc.)
	StaticDir fs.FS

	// StaticRouter provides a route that exposes the StaticDir
	StaticRouter struct {
		dir fs.FS
	}
)

// NewStaticRouter instantiates a new StaticRouter
func NewStaticRouter(staticDir StaticDir, appConfig cconfig.Config) (*StaticRouter, error) {
	var config config

	router := StaticRouter{
		dir: staticDir,
	}

	err := appConfig.Load("chttp", &config)
	if err != nil {
		return nil, cerrors.New(err, "failed to load chttp config", nil)
	}

	if config.WebDir != "" {
		router.dir = os.DirFS(config.WebDir)
	}

	return &router, nil
}

// Routes defines the /static/ route that serves the StaticDir
func (ro *StaticRouter) Routes() []Route {
	return []Route{
		{
			Path:    "/static/{path:.*}",
			Methods: []string{http.MethodGet},
			Handler: http.FileServer(http.FS(ro.dir)).ServeHTTP,
		},
	}
}
