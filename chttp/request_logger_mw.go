package chttp

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/tusharsoni/copper/clogger"
)

var errRWIsNotHijacker = errors.New("internal response writer is not http.Hijacker")

// NewRequestLoggerMiddleware creates a new RequestLoggerMiddleware.
func NewRequestLoggerMiddleware(logger clogger.Logger) *RequestLoggerMiddleware {
	return &RequestLoggerMiddleware{logger: logger}
}

// RequestLoggerMiddleware logs each request's HTTP method, path, and status code along with user uuid
// (from basic auth) if any.
type RequestLoggerMiddleware struct {
	logger clogger.Logger
}

// Handle wraps the current request with a request/response recorder. It records the method path and the
// return status code. It logs this with the given logger.
func (mw *RequestLoggerMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			tags     = make(map[string]interface{})
			loggerRw = requestLoggerRw{
				internal:   w,
				statusCode: http.StatusOK,
			}
		)

		user, _, ok := r.BasicAuth()
		if ok {
			tags["user"] = user
		}

		next.ServeHTTP(&loggerRw, r)

		mw.logger.WithTags(tags).Info(fmt.Sprintf("%s %s %d", r.Method, r.URL.Path, loggerRw.statusCode))
	})
}

type requestLoggerRw struct {
	internal   http.ResponseWriter
	statusCode int
}

func (rw *requestLoggerRw) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := rw.internal.(http.Hijacker)
	if !ok {
		return nil, nil, errRWIsNotHijacker
	}

	return h.Hijack()
}

func (rw *requestLoggerRw) Header() http.Header {
	return rw.internal.Header()
}

func (rw *requestLoggerRw) Write(b []byte) (int, error) {
	return rw.internal.Write(b)
}

func (rw *requestLoggerRw) WriteHeader(statusCode int) {
	rw.internal.WriteHeader(statusCode)
	rw.statusCode = statusCode
}
