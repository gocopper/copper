package chttp

import (
	"encoding/json"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/tusharsoni/copper/cerror"

	"github.com/tusharsoni/copper/clogger"
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

type responder struct {
	logger clogger.Logger
}

func NewJSONReaderWriter(logger clogger.Logger) ReaderWriter {
	govalidator.SetFieldsRequiredByDefault(true)

	return &responder{
		logger: logger,
	}
}

func (r *responder) Read(w http.ResponseWriter, req *http.Request, body interface{}) bool {
	url := req.URL.String()

	err := json.NewDecoder(req.Body).Decode(body)
	if err != nil {
		r.logger.Warn("Failed to read body", cerror.New(err, "invalid json", map[string]interface{}{
			"url": url,
		}))

		r.BadRequest(w, err)
		return false
	}

	ok, err := govalidator.ValidateStruct(body)
	if !ok {
		r.logger.Warn("Failed to read body", cerror.New(err, "data validation failed", map[string]interface{}{
			"url": url,
		}))

		r.BadRequest(w, err)
		return false
	}

	return true
}

func (r *responder) OK(w http.ResponseWriter, o interface{}) {
	r.json(w, o, http.StatusOK)
}

func (r *responder) Created(w http.ResponseWriter, o interface{}) {
	r.json(w, o, http.StatusCreated)
}

func (*responder) InternalErr(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}

func (*responder) Unauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
}

func (*responder) Forbidden(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
}

func (r *responder) BadRequest(w http.ResponseWriter, err error) {
	var resp struct {
		Error string `json:"error"`
	}
	resp.Error = err.Error()
	r.json(w, &resp, http.StatusBadRequest)
}

func (r *responder) json(w http.ResponseWriter, o interface{}, status int) {
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
