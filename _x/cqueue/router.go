package cqueue

import (
	"net/http"

	"github.com/gorilla/mux"
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
		svc:    p.Svc,
	}
}

func (ro *Router) Routes() []chttp.Route {
	return []chttp.Route{
		NewGetTaskRoute(ro).Route,
	}
}

func NewGetTaskRoute(ro *Router) chttp.RouteResult {
	return chttp.RouteResult{Route: chttp.Route{
		Path:    "/api/queue/tasks/{uuid}",
		Methods: []string{http.MethodGet},
		Handler: http.HandlerFunc(ro.HandleGetTask),
	}}
}

func (ro *Router) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	taskUUID := mux.Vars(r)["uuid"]

	task, err := ro.svc.GetTask(r.Context(), taskUUID)
	if err != nil {
		ro.logger.Error("Failed to get task", err)
		ro.rw.InternalErr(w)
		return
	}

	ro.rw.OK(w, task)
}
