package chttp

import (
	"encoding/json"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/tusharsoni/copper/v2/cerrors"
	"github.com/tusharsoni/copper/v2/clogger"
)

// ReaderWriter provides methods to read request body and write an http response easily.
type ReaderWriter interface {
	Read(w http.ResponseWriter, r *http.Request, body interface{}) bool

	OK(w http.ResponseWriter, o interface{})
	Created(w http.ResponseWriter, o interface{})
	InternalErr(w http.ResponseWriter)
	Unauthorized(w http.ResponseWriter)
	Forbidden(w http.ResponseWriter)
	BadRequest(w http.ResponseWriter, err error)
}

type jsonRW struct {
	logger clogger.Logger
}

// NewJSONReaderWriter creates a ReaderWriter that reads & responds using the JSON format.
// It validates reads using github.com/asaskevich/govalidator.
// During writes, the appropriate content-type header is set.
func NewJSONReaderWriter(logger clogger.Logger) ReaderWriter {
	return &jsonRW{
		logger: logger,
	}
}

func (r *jsonRW) Read(w http.ResponseWriter, req *http.Request, body interface{}) bool {
	url := req.URL.String()

	err := json.NewDecoder(req.Body).Decode(body)
	if err != nil {
		r.logger.Warn("Failed to read body", cerrors.New(err, "invalid json", map[string]interface{}{
			"url": url,
		}))

		r.BadRequest(w, err)

		return false
	}

	ok, err := govalidator.ValidateStruct(body)
	if !ok {
		r.logger.Warn("Failed to read body", cerrors.New(err, "data validation failed", map[string]interface{}{
			"url": url,
		}))

		r.BadRequest(w, err)

		return false
	}

	return true
}

func (r *jsonRW) OK(w http.ResponseWriter, o interface{}) {
	r.json(w, o, http.StatusOK)
}

func (r *jsonRW) Created(w http.ResponseWriter, o interface{}) {
	r.json(w, o, http.StatusCreated)
}

func (*jsonRW) InternalErr(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}

func (*jsonRW) Unauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
}

func (*jsonRW) Forbidden(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
}

func (r *jsonRW) BadRequest(w http.ResponseWriter, err error) {
	var resp struct {
		Error string `json:"error"`
	}

	resp.Error = err.Error()
	r.json(w, &resp, http.StatusBadRequest)
}

func (r *jsonRW) json(w http.ResponseWriter, o interface{}, status int) {
	j, err := json.Marshal(o)
	if err != nil {
		r.logger.Error("Failed to marshal response as json", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(j)
	if err != nil {
		r.logger.Error("Failed to write response to body", err)
	}
}
