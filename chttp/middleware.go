package chttp

import "net/http"

// Middleware defines the interface with a Handle func which is of the MiddlewareFunc type. Implementations of
// Middleware can be used with NewHandler for global middlewares or Route for route-specific middlewares.
type Middleware interface {
	Handle(next http.Handler) http.Handler
}

// MiddlewareFunc is a function that takes in a http.Handler and returns one as well. It allows you to execute
// code before or after calling the handler.
type MiddlewareFunc func(next http.Handler) http.Handler

// HandleMiddleware returns an implementation of Middleware that runs the provided func.
func HandleMiddleware(fn MiddlewareFunc) Middleware {
	return &middlewareFuncHandler{fn: fn}
}

type middlewareFuncHandler struct {
	fn MiddlewareFunc
}

func (mw *middlewareFuncHandler) Handle(next http.Handler) http.Handler {
	return mw.fn(next)
}
