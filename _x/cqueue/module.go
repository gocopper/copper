package cqueue

import (
	"github.com/tusharsoni/copper/chttp"

	"github.com/tusharsoni/copper/clogger"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Module struct {
	Svc    Svc
	Router *Router
}

type NewModuleParams struct {
	DB     *gorm.DB
	Config Config
	Logger clogger.Logger
}

func NewModule(p NewModuleParams) *Module {
	svc := NewSvc(SvcParams{
		Repo:   NewSQLRepo(p.DB),
		Config: p.Config,
	})

	return &Module{
		Svc: svc,
		Router: NewRouter(RouterParams{
			RW:     chttp.NewJSONReaderWriter(p.Logger),
			Logger: p.Logger,
			Svc:    svc,
		}),
	}
}

var Fx = fx.Provide(
	NewSQLRepo,
	NewSvc,

	NewRouter,
	NewGetTaskRoute,
)
