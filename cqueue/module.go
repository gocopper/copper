package cqueue

import (
	"github.com/tusharsoni/copper/clogger"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Module struct {
	Queue Svc
}

type NewModuleParams struct {
	DB      *gorm.DB
	Workers []Worker
	Config  Config
	Logger  clogger.Logger
}

func NewModule(p NewModuleParams) *Module {
	return &Module{NewSvc(SvcParams{
		Repo:    NewSQLRepo(p.DB),
		Workers: p.Workers,
		Config:  p.Config,
		Logger:  p.Logger,
	})}
}

var Fx = fx.Provide(
	NewSQLRepo,
	NewSvc,
)
