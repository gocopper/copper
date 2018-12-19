package chttp

import (
	"encoding/json"
	"net/http"

	"github.com/tusharsoni/copper/cerror"

	"github.com/tusharsoni/copper/clogger"

	"github.com/asaskevich/govalidator"
)

type BodyReader interface {
	Read(w http.ResponseWriter, r *http.Request, body interface{}) bool
}

type bodyReader struct {
	resp   Responder
	logger clogger.Logger
}

func newBodyReader(resp Responder, logger clogger.Logger) BodyReader {
	return &bodyReader{
		resp:   resp,
		logger: logger,
	}
}

func (b *bodyReader) Read(w http.ResponseWriter, r *http.Request, body interface{}) bool {
	url := r.URL.String()

	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		b.logger.Warn("Failed to read body", cerror.New(err, "invalid json", map[string]string{
			"url": url,
		}))

		b.resp.BadRequest(w, err)
		return false
	}

	govalidator.SetFieldsRequiredByDefault(true)

	ok, err := govalidator.ValidateStruct(body)
	if !ok {
		b.logger.Warn("Failed to read body", cerror.New(err, "data validation failed", map[string]string{
			"url": url,
		}))

		b.resp.BadRequest(w, err)
		return false
	}

	return true
}
