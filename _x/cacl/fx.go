package cacl

import (
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// Fx module for the cacl package that provides the SQL implementation for all services.
var Fx = fx.Provide(
	newSQLRepo,
	newSvcImpl,
)

type Module struct {
	Svc Svc
}

func NewModule(db *gorm.DB) *Module {
	return &Module{Svc: newSvcImpl(newSQLRepo(db))}
}
