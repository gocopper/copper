package anonymous

import (
	"net/http"

	"github.com/tusharsoni/copper/chttp"
	"github.com/tusharsoni/copper/clogger"
	"go.uber.org/fx"
)

type Router struct {
	rw     chttp.ReaderWriter
	logger clogger.Logger

	svc Svc
}

type RouterParams struct {
	fx.In

	RW     chttp.ReaderWriter
	Logger clogger.Logger

	Svc Svc
}

func NewRouter(p RouterParams) *Router {
	return &Router{
		rw:     p.RW,
		logger: p.Logger,

		svc: p.Svc,
	}
}

func NewCreateSessionRoute(ro *Router) chttp.RouteResult {
	return chttp.RouteResult{Route: chttp.Route{
		Path:    "/api/auth/anonymous/create",
		Methods: []string{http.MethodPost},
		Handler: http.HandlerFunc(ro.HandleCreateSession),
	}}
}

func (ro *Router) HandleCreateSession(w http.ResponseWriter, r *http.Request) {
	user, sessionToken, err := ro.svc.CreateAnonymousUser(r.Context())
	if err != nil {
		ro.logger.Error("Failed to create session token", err)
		ro.rw.InternalErr(w)
		return
	}

	ro.rw.OK(w, map[string]string{
		"user_uuid":     user.UUID,
		"session_token": sessionToken,
	})
}
