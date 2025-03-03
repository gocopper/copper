package chttp

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gocopper/copper/cerrors"
	"github.com/gocopper/copper/clogger"
)

type (
	// WriteJSONParams holds the params for the WriteJSON function in JSONReaderWriter
	WriteJSONParams struct {
		StatusCode int
		Data       interface{}
	}

	// JSONReaderWriter provides functions to read and write JSON data from/to HTTP requests/responses
	JSONReaderWriter struct {
		config Config
		logger clogger.Logger
	}
)

// NewJSONReaderWriter instantiates a new JSONReaderWriter with its dependencies
func NewJSONReaderWriter(config Config, logger clogger.Logger) *JSONReaderWriter {
	return &JSONReaderWriter{
		config: config,
		logger: logger,
	}
}

// WriteJSON writes a JSON response to the http.ResponseWriter. It can be configured with status code and data using
// WriteJSONParams.
func (rw *JSONReaderWriter) WriteJSON(w http.ResponseWriter, p WriteJSONParams) {
	w.Header().Set("Content-Type", "application/json")

	if p.StatusCode > 0 {
		w.WriteHeader(p.StatusCode)
	}

	if p.Data == nil {
		return
	}

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
func (rw *JSONReaderWriter) ReadJSON(w http.ResponseWriter, req *http.Request, body interface{}) bool {
	url := req.URL.String()

	err := json.NewDecoder(req.Body).Decode(body)
	if err != nil && errors.Is(err, io.EOF) {
		rw.WriteJSON(w, WriteJSONParams{
			StatusCode: http.StatusBadRequest,
			Data:       map[string]string{"error": "empty body"},
		})

		return false
	} else if err != nil {
		rw.WriteJSON(w, WriteJSONParams{
			StatusCode: http.StatusBadRequest,
			Data: cerrors.New(err, "invalid body json", map[string]interface{}{
				"url": url,
			}),
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

// Unauthorized writes a 401 Unauthorized JSON response
func (rw *JSONReaderWriter) Unauthorized(w http.ResponseWriter) {
	rw.WriteJSON(w, WriteJSONParams{
		StatusCode: http.StatusUnauthorized,
		Data:       map[string]string{"error": "unauthorized"},
	})
}
