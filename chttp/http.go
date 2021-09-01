package chttp

import (
	"net/http"
)

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
