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
	Routers           []Router
	GlobalMiddlewares []Middleware
}

// NewHandler creates a http.Handler with the given routes and middlewares.
// The handler can be used with a http.Server or as an argument to StartServer.
func NewHandler(p NewHandlerParams) http.Handler {
	var (
		muxRouter  = mux.NewRouter()
		muxHandler = http.NewServeMux()
	)

	for _, mw := range p.GlobalMiddlewares {
		muxRouter.Use(mw.Handle)
	}

	routes := make([]Route, 0)
	for _, router := range p.Routers {
		routes = append(routes, router.Routes()...)
	}

	sortRoutes(routes)

	for _, route := range routes {
		handler := http.Handler(route.Handler)

		// Since we are wrapping the handler in middleware functions, the outermost one will run first.
		// By applying the middlewares in reverse, we ensure that the first middleware in the list is the outermost one.
		for i := len(route.Middlewares) - 1; i >= 0; i-- {
			handler = route.Middlewares[i].Handle(handler)
		}

		handler = setRoutePathInCtxMiddleware(route.Path).Handle(handler)

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
