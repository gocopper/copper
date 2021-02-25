package chttptest

import (
	"net/http"
)

// NewMiddleware returns an implementation of chttp.Middleware that runs the provided func.
func NewMiddleware(fn func(next http.Handler) http.Handler) *Middleware {
	return &Middleware{fn: fn}
}

// Middleware is a simple implementation of chttp.Middleware useful for testing that simply
// runs the provided func.
type Middleware struct {
	fn func(next http.Handler) http.Handler
}

// Handle runs the provided func and returns its result.
func (mw *Middleware) Handle(next http.Handler) http.Handler {
	return mw.fn(next)
}
