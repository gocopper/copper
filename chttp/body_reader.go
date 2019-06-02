package chttp

import (
	"encoding/json"
	"net/http"

	"github.com/tusharsoni/copper/cerror"

	"github.com/tusharsoni/copper/clogger"

	"github.com/asaskevich/govalidator"
)

// BodyReader provides methods to read the incoming http request's JSON body, parse it into a struct, validate the
// data using govalidator, and respond if it's a bad request.
type BodyReader interface {
	Read(w http.ResponseWriter, r *http.Request, body interface{}) bool
}

type bodyReader struct {
	resp   Responder
	logger clogger.Logger
}

func newBodyReader(resp Responder, logger clogger.Logger) BodyReader {
	govalidator.SetFieldsRequiredByDefault(true)

	return &bodyReader{
		resp:   resp,
		logger: logger,
	}
}

func (b *bodyReader) Read(w http.ResponseWriter, r *http.Request, body interface{}) bool {
	url := r.URL.String()

	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		b.logger.Warn("Failed to read body", cerror.New(err, "invalid json", map[string]interface{}{
			"url": url,
		}))

		b.resp.BadRequest(w, err)
		return false
	}

	ok, err := govalidator.ValidateStruct(body)
	if !ok {
		b.logger.Warn("Failed to read body", cerror.New(err, "data validation failed", map[string]interface{}{
			"url": url,
		}))

		b.resp.BadRequest(w, err)
		return false
	}

	return true
}
