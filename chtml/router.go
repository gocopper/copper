package chtml

import (
	"net/http"
	"path"

	"github.com/gocopper/copper/chttp"
)

type (
	// Router provides routes to serve static assets and the index page for a single-page app
	Router struct {
		rw        *ReaderWriter
		staticDir StaticDir
		config    Config
	}

	// NewRouterParams holds the params needed to instantiate a new Router
	NewRouterParams struct {
		StaticDir StaticDir
		RW        *ReaderWriter
		Config    Config
	}
)

// NewRouter instantiates a new Router
func NewRouter(p NewRouterParams) (*Router, error) {
	return &Router{
		rw:        p.RW,
		staticDir: p.StaticDir,
		config:    p.Config,
	}, nil
}

// Routes defines the HTTP routes for this router
func (ro *Router) Routes() []chttp.Route {
	routes := []chttp.Route{
		{
			Path:    "/static/{path:.*}",
			Methods: []string{http.MethodGet},
			Handler: ro.HandleStaticFile,
		},
	}

	if ro.config.EnableSinglePageRouting {
		routes = append(routes, chttp.Route{
			Path:    "/{path:.*}",
			Methods: []string{http.MethodGet},
			Handler: ro.HandleIndexPage,
		})
	}

	return routes
}

// HandleStaticFile serves the requested static file as found in the web/public directory. In non-dev env, the static
// files are embedded in the binary.
func (ro *Router) HandleStaticFile(w http.ResponseWriter, r *http.Request) {
	if ro.config.DevMode {
		http.ServeFile(w, r, path.Join("web", "public", chttp.URLParams(r)["path"]))
		return
	}

	http.FileServer(http.FS(ro.staticDir)).ServeHTTP(w, r)
}

// HandleIndexPage renders the index.html page
func (ro *Router) HandleIndexPage(w http.ResponseWriter, r *http.Request) {
	ro.rw.WriteHTML(w, r, WriteHTMLParams{
		PageTemplate: "index.html",
	})
}
