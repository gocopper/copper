package chttp

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/asaskevich/govalidator"
	"github.com/gocopper/copper/cerrors"
	"github.com/gocopper/copper/clogger"
)

type (
	// WriteJSONParams holds the params for the WriteJSON function in ReaderWriter
	WriteJSONParams struct {
		StatusCode int
		Data       interface{}
	}

	// ReaderWriter provides functions to read data from HTTP requests and write response bodies in various formats
	ReaderWriter struct {
		logger clogger.Logger
	}
)

// URLParams returns the route variables for the current request, if any
var URLParams = mux.Vars

// NewReaderWriter instantiates a new ReaderWriter with its dependencies
func NewReaderWriter(logger clogger.Logger) *ReaderWriter {
	return &ReaderWriter{
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
