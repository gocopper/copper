package chtml

import (
	"net/http"

	"github.com/gocopper/copper/chttp"
)

type (
	// Router provides routes to serve static assets and the index page for a single-page app
	Router struct {
		rw     *ReaderWriter
		config Config

		staticFileServer http.Handler
	}

	// NewRouterParams holds the params needed to instantiate a new Router
	NewRouterParams struct {
		StaticDir StaticDir
		RW        *ReaderWriter
		Config    Config
	}
)

// NewRouter instantiates a new Router
func NewRouter(p NewRouterParams) *Router {
	return &Router{
		rw:     p.RW,
		config: p.Config,

		staticFileServer: http.FileServer(http.FS(p.StaticDir)),
	}
}

// Routes defines the HTTP routes for this router
func (ro *Router) Routes() []chttp.Route {
	routes := []chttp.Route{
		{
			Path:    "/static/{path:.*}",
			Methods: []string{http.MethodGet},
			Handler: ro.staticFileServer.ServeHTTP,
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

// HandleIndexPage renders the index.html page
func (ro *Router) HandleIndexPage(w http.ResponseWriter, r *http.Request) {
	ro.rw.WriteHTML(w, r, WriteHTMLParams{
		PageTemplate: "index.html",
	})
}
