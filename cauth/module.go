package cauth

import (
	"github.com/jinzhu/gorm"
	"github.com/tusharsoni/copper/chttp"
	"github.com/tusharsoni/copper/clogger"
)

type Module struct {
	Svc        Svc
	Middleware Middleware
}

func NewModule(db *gorm.DB, logger clogger.Logger) *Module {
	var (
		svc = NewSvc(NewSQLRepo(db))
		mw  = NewAuthMiddleware(MiddlewareParams{
			RW:     chttp.NewJSONReaderWriter(logger),
			Svc:    svc,
			Logger: logger,
		})
	)

	return &Module{
		Svc:        svc,
		Middleware: mw,
	}
}
