package cqueue

import (
	"context"

	"github.com/tusharsoni/copper/clogger"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Module struct {
	Svc Svc
}

type NewModuleParams struct {
	DB     *gorm.DB
	Config Config
}

func NewModule(p NewModuleParams) *Module {
	return &Module{NewSvc(SvcParams{
		Repo:   NewSQLRepo(p.DB),
		Config: p.Config,
	})}
}

var Fx = fx.Provide(
	NewSQLRepo,
	NewSvc,
)

type StartBackgroundWorkersParams struct {
	fx.In

	Workers []Worker `group:"cqueue/workers"`
	Queue   Svc
	Logger  clogger.Logger
}

func StartBackgroundWorkers(p StartBackgroundWorkersParams) {
	for i := range p.Workers {
		go p.Workers[i].Start(context.Background(), p.Queue, p.Logger)
	}
}
