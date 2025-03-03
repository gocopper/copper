package chttp

import (
	// Used to embed error.html
	_ "embed"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/gocopper/copper/cerrors"
	"github.com/gocopper/copper/clogger"
)

type (
	// WriteHTMLParams holds the params for the WriteHTML function in ReaderWriter
	WriteHTMLParams struct {
		StatusCode     int
		Error          error
		Data           interface{}
		PageTemplate   string
		LayoutTemplate string
	}

	// WritePartialParams are the parameters for WritePartial
	WritePartialParams struct {
		Name string
		Data interface{}
	}

	// HTMLReaderWriter provides functions to read data from HTTP requests and write HTML response bodies
	HTMLReaderWriter struct {
		html   *HTMLRenderer
		config Config
		logger clogger.Logger
	}
)

// URLParams returns the route variables for the current request, if any
var URLParams = mux.Vars

//go:embed error.html
var errorHTML string

// NewHTMLReaderWriter instantiates a new HTMLReaderWriter with its dependencies
func NewHTMLReaderWriter(html *HTMLRenderer, config Config, logger clogger.Logger) *HTMLReaderWriter {
	return &HTMLReaderWriter{
		html:   html,
		config: config,
		logger: logger,
	}
}

// WriteHTMLError handles the given error. In render_error is configured to true, it writes an HTML page with the error.
// Errors are always logged.
func (rw *HTMLReaderWriter) WriteHTMLError(w http.ResponseWriter, r *http.Request, err error) {
	rw.WriteHTML(w, r, WriteHTMLParams{
		Error: err,
	})
}

// WriteHTML writes an HTML response to the provided http.ResponseWriter. Using the given WriteHTMLParams, the HTML
// is generated with a layout, page, and component templates.
func (rw *HTMLReaderWriter) WriteHTML(w http.ResponseWriter, r *http.Request, p WriteHTMLParams) {
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

	if p.Error != nil && rw.config.RenderHTMLError {
		w.Header().Set("content-type", "text/html")
		w.WriteHeader(p.StatusCode)

		errorHTMLTmpl := template.Must(template.New("chtml/error.html").Parse(errorHTML))

		_ = errorHTMLTmpl.Execute(w, map[string]interface{}{
			"Error": p.Error.Error(),
		})

		return
	}

	out, err := rw.html.render(r, p.LayoutTemplate, p.PageTemplate, p.Data)
	if err != nil {
		rw.logger.Error("Failed to render html template", cerrors.WithTags(err, map[string]interface{}{
			"layout": p.LayoutTemplate,
			"page":   p.PageTemplate,
		}))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/html")
	w.WriteHeader(p.StatusCode)
	_, _ = w.Write([]byte(out))
}

// WritePartial renders a partial template with the given name and data
func (rw *HTMLReaderWriter) WritePartial(w http.ResponseWriter, r *http.Request, p WritePartialParams) {
	out, err := rw.html.partial(r)(p.Name, p.Data)
	if err != nil {
		rw.WriteHTMLError(w, r, cerrors.New(err, "failed to render partial", map[string]interface{}{
			"name": p.Name,
		}))
		return
	}

	w.Header().Set("content-type", "text/html")
	_, _ = w.Write([]byte(out))
}

// Unauthorized writes a 401 Unauthorized response to the http.ResponseWriter. If a redirect URL is configured,
// the user is redirected to that URL instead.
func (rw *HTMLReaderWriter) Unauthorized(w http.ResponseWriter, r *http.Request) {
	if rw.config.RedirectURLForUnauthorizedRequests != nil {
		http.Redirect(w, r, *rw.config.RedirectURLForUnauthorizedRequests, http.StatusSeeOther)
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
}
