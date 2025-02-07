package chttp

import (
	"net/http"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/gocopper/copper/clogger"

	"github.com/gorilla/mux"
)

// NewHandlerParams holds the params needed for NewHandler.
type NewHandlerParams struct {
	Routers           []Router
	GlobalMiddlewares []Middleware
	Config            Config
	Logger            clogger.Logger
}

// NewHandler creates a http.Handler with the given routes and middlewares.
// The handler can be used with a http.Server or as an argument to StartServer.
func NewHandler(p NewHandlerParams) http.Handler {
	var (
		muxRouter  = mux.NewRouter()
		muxHandler = http.NewServeMux()
	)

	routes := make([]Route, 0)
	for _, router := range p.Routers {
		routes = append(routes, router.Routes()...)
	}

	sortRoutes(routes)

	for _, route := range routes {
		routePath := route.Path

		// If a base path is set in the configuration ensure that the route path
		// starts with the base path. If it does not, skip this route.
		// Otherwise, remove the base path prefix from the route path to allow for
		// proper registration in the router.
		if p.Config.BasePath != nil {
			if route.RegisterWithBasePath {
				routePath = path.Join(*p.Config.BasePath, routePath)
			}

			if !strings.HasPrefix(routePath, *p.Config.BasePath) {
				continue
			}

			routePath = routePath[len(*p.Config.BasePath):]
		}

		handler := http.Handler(route.Handler)

		// Register route-level handlers
		// Since we are wrapping the handler in middleware functions, the outermost one will run first.
		// By applying the middlewares in reverse, we ensure that the first middleware in the list is the outermost one.
		for i := len(route.Middlewares) - 1; i >= 0; i-- {
			handler = route.Middlewares[i].Handle(handler)
		}

		// Register global middlewares
		for i := len(p.GlobalMiddlewares) - 1; i >= 0; i-- {
			handler = p.GlobalMiddlewares[i].Handle(handler)
		}

		handler = setRoutePathInCtxMiddleware(routePath).Handle(handler)
		handler = panicLoggerMiddleware(p.Logger).Handle(handler)

		muxRoute := muxRouter.Handle(routePath, handler)

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
