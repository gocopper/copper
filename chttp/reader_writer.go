package chttp

import (
	// Used to embed error.html
	_ "embed"
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/asaskevich/govalidator"
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

	// WriteJSONParams holds the params for the WriteJSON function in ReaderWriter
	WriteJSONParams struct {
		StatusCode int
		Data       interface{}
	}

	// ReaderWriter provides functions to read data from HTTP requests and write response bodies in various formats
	ReaderWriter struct {
		html   *HTMLRenderer
		config Config
		logger clogger.Logger
	}
)

// URLParams returns the route variables for the current request, if any
var URLParams = mux.Vars

//go:embed error.html
var errorHTML string

// NewReaderWriter instantiates a new ReaderWriter with its dependencies
func NewReaderWriter(html *HTMLRenderer, config Config, logger clogger.Logger) *ReaderWriter {
	return &ReaderWriter{
		html:   html,
		config: config,
		logger: logger,
	}
}

// WriteJSON writes a JSON response to the http.ResponseWriter. It can be configured with status code and data using
// WriteJSONParams.
func (rw *ReaderWriter) WriteJSON(w http.ResponseWriter, p WriteJSONParams) {
	if p.StatusCode > 0 {
		w.WriteHeader(p.StatusCode)
	}

	if p.Data == nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	errData, ok := p.Data.(error)
	if ok {
		err := json.NewEncoder(w).Encode(map[string]string{
			"error": errData.Error(),
		})
		if err != nil {
			rw.logger.Error("Failed to marshal error response as json", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	err := json.NewEncoder(w).Encode(p.Data)
	if err != nil {
		rw.logger.Error("Failed to marshal response as json", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

// ReadJSON reads JSON from the http.Request into the body var. If the body struct has validate tags on it, the
// struct is also validated. If the validation fails, a BadRequest response is sent back and the function returns
// false.
func (rw *ReaderWriter) ReadJSON(w http.ResponseWriter, req *http.Request, body interface{}) bool {
	url := req.URL.String()

	err := json.NewDecoder(req.Body).Decode(body)
	if err != nil {
		rw.logger.Warn("Failed to read body", cerrors.New(err, "invalid json", map[string]interface{}{
			"url": url,
		}))

		rw.WriteJSON(w, WriteJSONParams{
			StatusCode: http.StatusBadRequest,
			Data:       err,
		})

		return false
	}

	ok, err := govalidator.ValidateStruct(body)
	if !ok {
		rw.logger.Warn("Failed to read body", cerrors.New(err, "data validation failed", map[string]interface{}{
			"url": url,
		}))

		rw.WriteJSON(w, WriteJSONParams{
			StatusCode: http.StatusBadRequest,
			Data:       err,
		})

		return false
	}

	return true
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
	var renderHTMLError = rw.config.DevMode

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

	if p.Error != nil && renderHTMLError {
		w.WriteHeader(p.StatusCode)
		w.Header().Set("content-type", "text/html")

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

	w.WriteHeader(p.StatusCode)
	w.Header().Set("content-type", "text/html")
	_, _ = w.Write([]byte(out))
}
