package chttp

import (
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/gorilla/mux"
)

// NewHandlerParams holds the params needed for NewHandler.
type NewHandlerParams struct {
	Routes            []Route
	GlobalMiddlewares []Middleware
}

// NewHandler creates a http.Handler with the given routes and middlewares.
// The handler can be used with a http.Server or as an argument to StartServer.
func NewHandler(p NewHandlerParams) http.Handler {
	var (
		muxRouter  = mux.NewRouter()
		muxHandler = http.NewServeMux()
	)

	for _, f := range p.GlobalMiddlewares {
		muxRouter.Use(mux.MiddlewareFunc(f))
	}

	sortRoutes(p.Routes)

	for _, route := range p.Routes {
		handler := http.Handler(route.Handler)

		for _, f := range route.Middlewares {
			handler = f(handler)
		}

		muxRoute := muxRouter.Handle(route.Path, handler)

		if len(route.Methods) > 0 {
			muxRoute.Methods(route.Methods...)
		}
	}

	muxHandler.Handle("/", muxRouter)

	return muxHandler
}

func sortRoutes(routes []Route) {
	const matcherPlaceholder = "{{matcher}}"

	re := regexp.MustCompile(`(?U)(\{.*\})`)

	sort.Slice(routes, func(i, j int) bool {
		aPath := re.ReplaceAllString(routes[i].Path, matcherPlaceholder)
		bPath := re.ReplaceAllString(routes[j].Path, matcherPlaceholder)

		aParts := strings.Split(aPath, "/")
		bParts := strings.Split(bPath, "/")

		if aPath == "/" {
			aParts = nil
		}

		if bPath == "/" {
			bParts = nil
		}

		if len(aParts) != len(bParts) {
			return len(aParts) > len(bParts)
		}

		for i, aPart := range aParts {
			bPart := bParts[i]

			if aPart == matcherPlaceholder {
				return false
			}

			if bPart == matcherPlaceholder {
				return true
			}
		}

		return false
	})
}
