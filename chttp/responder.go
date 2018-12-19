package chttp

import (
	"encoding/json"
	"net/http"

	"github.com/tusharsoni/copper/clogger"
)

type Responder struct {
	logger clogger.Logger
}

func newResponder(logger clogger.Logger) *Responder {
	return &Responder{
		logger: logger,
	}
}

func (r *Responder) OK(w http.ResponseWriter, o interface{}) {
	r.json(w, o, http.StatusOK)
}

func (r *Responder) Created(w http.ResponseWriter, o interface{}) {
	r.json(w, o, http.StatusCreated)
}

func (*Responder) InternalErr(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}

func (*Responder) Unauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
}

func (*Responder) Forbidden(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
}

func (r *Responder) BadRequest(w http.ResponseWriter, err error) {
	var resp struct {
		Error string `json:"error"`
	}
	resp.Error = err.Error()
	r.json(w, &resp, http.StatusBadRequest)
}

func (r *Responder) json(w http.ResponseWriter, o interface{}, status int) {
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
