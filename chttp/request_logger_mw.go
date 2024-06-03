package chttp

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/gocopper/copper/cmetrics"
	"net"
	"net/http"
	"time"

	"github.com/gocopper/copper/clogger"
)

var errRWIsNotHijacker = errors.New("internal response writer is not http.Hijacker")

// NewRequestLoggerMiddleware creates a new RequestLoggerMiddleware.
func NewRequestLoggerMiddleware(metrics cmetrics.Metrics, logger clogger.Logger) *RequestLoggerMiddleware {
	return &RequestLoggerMiddleware{
		metrics: metrics,
		logger:  logger,
	}
}

// RequestLoggerMiddleware logs each request's HTTP method, path, and status code along with user uuid
// (from basic auth) if any.
type RequestLoggerMiddleware struct {
	metrics cmetrics.Metrics
	logger  clogger.Logger
}

// Handle wraps the current request with a request/response recorder. It records the method path and the
// return status code. It logs this with the given logger.
func (mw *RequestLoggerMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			loggerRw = requestLoggerRw{
				internal:   w,
				statusCode: http.StatusOK,
			}

			tags = map[string]interface{}{
				"method": r.Method,
				"url":    r.URL.Path,
			}

			begin = time.Now()
		)

		user, _, ok := r.BasicAuth()
		if ok {
			tags["user"] = user
		}

		next.ServeHTTP(&loggerRw, r)

		tags["statusCode"] = loggerRw.statusCode
		tags["duration"] = time.Since(begin).String()

		mw.metrics.CounterInc("http_requests_total", map[string]string{
			"status_code": fmt.Sprintf("%d", loggerRw.statusCode),
			"path":        RawRoutePath(r),
		})
		mw.metrics.HistogramObserve("http_request_duration_seconds", map[string]string{
			"status_code": fmt.Sprintf("%d", loggerRw.statusCode),
			"path":        RawRoutePath(r),
		}, time.Since(begin).Seconds())

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
