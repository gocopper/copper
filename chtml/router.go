package chtml

import (
	"io/fs"
	"net/http"

	"github.com/gocopper/copper/chttp"
)

type (
	// Router provides routes to serve static assets and the index page for a single-page app
	Router struct {
		dir fs.FS
		rw  *ReaderWriter
	}

	// NewRouterParams holds the params needed to instantiate a new Router
	NewRouterParams struct {
		StaticDir StaticDir
		RW        *ReaderWriter
	}
)

// NewRouter instantiates a new Router
func NewRouter(p NewRouterParams) *Router {
	return &Router{
		dir: p.StaticDir,
		rw:  p.RW,
	}
}

// Routes defines the HTTP routes for this router
func (ro *Router) Routes() []chttp.Route {
	return []chttp.Route{
		{
			Path:    "/{path:.*}",
			Methods: []string{http.MethodGet},
			Handler: ro.HandleIndexPage,
		},
		{
			Path:    "/static/{path:.*}",
			Methods: []string{http.MethodGet},
			Handler: http.FileServer(http.FS(ro.dir)).ServeHTTP,
		},
	}
}

// HandleIndexPage renders the index.html page
func (ro *Router) HandleIndexPage(w http.ResponseWriter, r *http.Request) {
	ro.rw.WriteHTML(w, r, WriteHTMLParams{
		PageTemplate: "index.html",
	})
}
