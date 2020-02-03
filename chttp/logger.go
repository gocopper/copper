package chttp

import (
	"fmt"
	"net/http"

	"github.com/tusharsoni/copper/clogger"
)

func newRequestLogger(logger clogger.Logger) GlobalMiddlewareFuncResult {
	var mw = func(next http.Handler) http.Handler {
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

			logger.WithTags(tags).Info(fmt.Sprintf("%s %s %d", r.Method, r.URL.Path, loggerRw.statusCode))
		})
	}

	return GlobalMiddlewareFuncResult{
		GlobalMiddlewareFunc: mw,
	}
}

type requestLoggerRw struct {
	internal   http.ResponseWriter
	statusCode int
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
