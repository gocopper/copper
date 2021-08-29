package chtml

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"

	// embeds html templates
	_ "embed"

	"github.com/gocopper/copper/cerrors"
	"github.com/gocopper/copper/clogger"
)

//go:embed error.html
var errorHTML string

type (
	// WriteHTMLParams holds the params for the WriteHTML function in ReaderWriter
	WriteHTMLParams struct {
		StatusCode     int
		Error          error
		Data           interface{}
		PageTemplate   string
		LayoutTemplate string
	}

	// ReaderWriter provides functions to read data from HTTP requests and write response bodies in various formats
	ReaderWriter struct {
		renderer *Renderer
		logger   clogger.Logger

		renderError bool
	}
)

// URLParams returns the route variables for the current request, if any
var URLParams = mux.Vars

// NewReaderWriter instantiates a new ReaderWriter with its dependencies
func NewReaderWriter(renderer *Renderer, config Config, logger clogger.Logger) *ReaderWriter {
	return &ReaderWriter{
		renderer: renderer,
		logger:   logger,

		renderError: config.DevMode,
	}
}

// WriteHTMLError handles the given error. In render_error is configured to true, it writes an HTML page with the error.
// Errors are always logged.
func (rw *ReaderWriter) WriteHTMLError(w http.ResponseWriter, r *http.Request, err error) {
	rw.WriteHTML(w, r, WriteHTMLParams{
		Error: err,
	})
}

// WriteHTML writes an HTML response to the provided http.ResponseWriter. Using the given WriteHTMLParams, the HTML
// is generated with a layout, page, and component templates.
func (rw *ReaderWriter) WriteHTML(w http.ResponseWriter, r *http.Request, p WriteHTMLParams) {
	if p.StatusCode == 0 && p.Error == nil {
		p.StatusCode = http.StatusOK
	}

	if p.StatusCode == 0 && p.Error != nil {
		p.StatusCode = http.StatusInternalServerError
	}

	if p.LayoutTemplate == "" {
		p.LayoutTemplate = "main.html"
	}

	if p.PageTemplate == "" && p.StatusCode == http.StatusInternalServerError {
		p.PageTemplate = "internal-error.html"
	}

	if p.PageTemplate == "" && p.StatusCode == http.StatusNotFound {
		p.PageTemplate = "not-found.html"
	}

	if p.Error != nil {
		rw.logger.WithTags(map[string]interface{}{
			"url": r.URL.String(),
		}).Error("Failed to handle request", p.Error)
	}

	if p.Error != nil && rw.renderError {
		w.WriteHeader(p.StatusCode)
		w.Header().Set("content-type", "text/html")

		errorHTMLTmpl := template.Must(template.New("chtml/error.html").Parse(errorHTML))

		_ = errorHTMLTmpl.Execute(w, map[string]interface{}{
			"Error": p.Error.Error(),
		})

		return
	}

	out, err := rw.renderer.render(r, p.LayoutTemplate, p.PageTemplate, p.Data)
	if err != nil {
		rw.logger.Error("Failed to render html template", cerrors.WithTags(err, map[string]interface{}{
			"layout": p.LayoutTemplate,
			"page":   p.PageTemplate,
		}))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(p.StatusCode)
	w.Header().Set("content-type", "text/html")
	_, _ = w.Write([]byte(out))
}
