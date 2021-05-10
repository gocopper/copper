package chttp

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"os"

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

	if config.WebDir != "" {
		router.dir = os.DirFS(config.WebDir)
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
		{
			Path:    "/api/chtml/components/call-method",
			Methods: []string{http.MethodPost},
			Handler: ro.HandleCallComponentMethod,
		},
	}
}

// HandleCallComponentMethod calls the action method on the component along with its props and args.
// The component is then re-rendered and the updated HTML along with any broadcasted events are
// sent back in the response.
func (ro *HTMLRouter) HandleCallComponentMethod(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID        string            `json:"id"`
		Component string            `json:"component"`
		Method    string            `json:"method"`
		Props     []json.RawMessage `json:"props"`
		Args      []json.RawMessage `json:"args"`
	}

	if !ro.rw.ReadJSON(w, r, &body) {
		return
	}

	req := requestWithComponentTree(r)

	html, err := ro.html.callComponentMethod(req, body.ID, body.Component, body.Method, body.Props, body.Args)
	if err != nil {
		ro.logger.Error("Failed to call component method", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ro.rw.WriteJSON(w, WriteJSONParams{
		Data: map[string]interface{}{
			"events": GetComponentTree(req).events,
			"html":   html,
		},
	})
}
