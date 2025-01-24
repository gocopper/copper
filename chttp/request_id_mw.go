package chttp

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ctxRequestID string

const ctxRequestIDKey = ctxRequestID("chttp/request-id")

// SetRequestIDInCtxMiddleware sets a unique request id in the context
func SetRequestIDInCtxMiddleware() Middleware {
	var mw = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), ctxRequestIDKey, uuid.New().String())

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	return HandleMiddleware(mw)
}

// GetRequestID returns the request id from the context.
// If the request id is not found in the context, it returns
// an empty string.
// It should be used only after the request id middleware is applied.
func GetRequestID(ctx context.Context) string {
	requestID, ok := ctx.Value(ctxRequestIDKey).(string)
	if !ok || requestID == "" {
		return ""
	}

	return requestID
}
