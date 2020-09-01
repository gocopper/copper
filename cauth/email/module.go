package email

import (
	"github.com/jinzhu/gorm"
	"github.com/tusharsoni/copper/cauth"
	"github.com/tusharsoni/copper/chttp"
	"github.com/tusharsoni/copper/clogger"
	"github.com/tusharsoni/copper/cmailer"
)

type Module struct {
	Svc    Svc
	Routes []chttp.Route
}

type NewModuleParams struct {
	Auth   cauth.Svc
	Mailer cmailer.Mailer
	Logger clogger.Logger
	Config Config

	DB *gorm.DB
}

func NewModule(p NewModuleParams) *Module {
	var (
		svc = NewSvc(SvcParams{
			Auth:   p.Auth,
			Repo:   NewSQLRepo(p.DB),
			Mailer: p.Mailer,
			Config: p.Config,
			Logger: p.Logger,
		})

		router = NewRouter(RouterParams{
			RW:     chttp.NewJSONReaderWriter(p.Logger),
			Logger: p.Logger,
			Auth:   svc,
			AuthMW: cauth.NewAuthMiddleware(cauth.MiddlewareParams{
				RW:     chttp.NewJSONReaderWriter(p.Logger),
				Svc:    p.Auth,
				Logger: p.Logger,
			}),
			Config: p.Config,
		})
	)

	return &Module{
		Svc:    svc,
		Routes: router.Routes(),
	}
}
