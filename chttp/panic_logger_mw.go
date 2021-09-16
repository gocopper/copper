package chttp

import (
	"net/http"

	"github.com/gocopper/copper/clogger"
)

func panicLoggerMiddleware(logger clogger.Logger) Middleware {
	mw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				log := logger.WithTags(map[string]interface{}{
					"path": r.URL.Path,
				})

				switch r := recover().(type) {
				case error:
					log.Error("Recovered from a panic while handling HTTP request", r)
				default:
					log.WithTags(map[string]interface{}{
						"error": r,
					}).Error("Recovered from a panic while handling HTTP request", nil)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}

	return HandleMiddleware(mw)
}
