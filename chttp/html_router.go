package chttp

import (
	"net/http"
	"path"
	"strings"
)

type (
	// HTMLRouter provides routes to serve (1) static assets (2) index page for an SPA
	HTMLRouter struct {
		rw        *HTMLReaderWriter
		staticDir StaticDir
		config    Config
	}

	// NewHTMLRouterParams holds the params needed to instantiate a new Router
	NewHTMLRouterParams struct {
		StaticDir StaticDir
		RW        *HTMLReaderWriter
		Config    Config
	}
)

// NewHTMLRouter instantiates a new Router
func NewHTMLRouter(p NewHTMLRouterParams) (*HTMLRouter, error) {
	return &HTMLRouter{
		rw:        p.RW,
		staticDir: p.StaticDir,
		config:    p.Config,
	}, nil
}

// Routes defines the HTTP routes for this router
func (ro *HTMLRouter) Routes() []Route {
	routes := []Route{
		{
			Path:                 "/static/{path:.*}",
			Methods:              []string{http.MethodGet},
			Handler:              ro.HandleStaticFile,
			RegisterWithBasePath: true,
		},
	}

	if ro.config.EnableSinglePageRouting {
		routes = append(routes, Route{
			Path:                 "/{path:.*}",
			Methods:              []string{http.MethodGet},
			Handler:              ro.HandleIndexPage,
			RegisterWithBasePath: true,
		})
	}

	return routes
}

// HandleStaticFile serves the requested static file as found in the web/public directory. In non-dev env, the static
// files are embedded in the binary.
func (ro *HTMLRouter) HandleStaticFile(w http.ResponseWriter, r *http.Request) {
	// Disable directory listing
	if strings.HasSuffix(r.URL.Path, "/") {
		http.NotFound(w, r)
		return
	}

	if ro.config.UseLocalHTML {
		http.ServeFile(w, r, path.Join("web", "public", URLParams(r)["path"]))
		return
	}

	http.FileServer(http.FS(ro.staticDir)).ServeHTTP(w, r)
}

// HandleIndexPage renders the index.html page
func (ro *HTMLRouter) HandleIndexPage(w http.ResponseWriter, r *http.Request) {
	ro.rw.WriteHTML(w, r, WriteHTMLParams{
		PageTemplate: "index.html",
	})
}
