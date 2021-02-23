package chttp

import (
	"net/http"
)

// Middleware is a function that takes in a http.Handler and returns one as well. It allows you to execute
// code before or after calling the handler.
type Middleware func(http.Handler) http.Handler

// Route represents a single HTTP route (ex. /api/profile) that can be configured with middlewares, path,
// HTTP methods, and a handler.
type Route struct {
	Middlewares []Middleware
	Path        string
	Methods     []string
	Handler     http.HandlerFunc
}

// Router is used to group routes together that are returned by the Routes method.
type Router interface {
	Routes() []Route
}
