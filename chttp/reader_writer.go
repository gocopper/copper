package chttp

import (
	"encoding/json"
	"html/template"
	"io/fs"
	"net/http"
	"path"

	"github.com/asaskevich/govalidator"
	"github.com/gocopper/copper/cerrors"
	"github.com/gocopper/copper/clogger"
)

type (
	// HTMLDir is a directory that can be embedded or found on the host system. It should contain sub-directories
	// and files to support the WriteHTML function in ReaderWriter.
	HTMLDir fs.FS

	// WriteJSONParams holds the params for the WriteJSON function in ReaderWriter
	WriteJSONParams struct {
		StatusCode int
		Data       interface{}
	}

	// WriteHTMLParams holds the params for the WriteHTML function in ReaderWriter
	WriteHTMLParams struct {
		StatusCode     int
		Data           interface{}
		PageTemplate   string
		LayoutTemplate string
	}

	// ReaderWriter provides functions to read data from HTTP requests and write response bodies in various formats
	ReaderWriter struct {
		htmlDir fs.FS
		logger  clogger.Logger
	}
)

// NewReaderWriter instantiates a new ReaderWriter with its dependencies
func NewReaderWriter(htmlDir HTMLDir, logger clogger.Logger) *ReaderWriter {
	return &ReaderWriter{
		htmlDir: htmlDir,
		logger:  logger,
	}
}

// WriteHTML writes an HTML response to the provided http.ResponseWriter. Using the given WriteHTMLParams, the HTML
// is generated with a layout, page, and component templates.
func (rw *ReaderWriter) WriteHTML(w http.ResponseWriter, p WriteHTMLParams) {
	if p.StatusCode == 0 {
		p.StatusCode = http.StatusOK
	}

	if p.LayoutTemplate == "" {
		p.LayoutTemplate = "main.html"
	}

	if p.PageTemplate == "" && p.StatusCode == http.StatusNotFound {
		p.PageTemplate = "not-found.html"
	}

	tmpl, err := template.ParseFS(rw.htmlDir,
		path.Join("html", "layouts", p.LayoutTemplate),
		path.Join("html", "pages", p.PageTemplate),
		path.Join("html", "components", "*.html"),
	)
	if err != nil {
		rw.logger.Error("Failed to render view template", cerrors.WithTags(err, map[string]interface{}{
			"layout": p.LayoutTemplate,
			"page":   p.PageTemplate,
		}))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(p.StatusCode)
	w.Header().Set("content-type", "text/html")

	err = tmpl.Execute(w, p.Data)
	if err != nil {
		rw.logger.Error("Failed to render view template", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
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
