package emailotp

import (
	"github.com/tusharsoni/copper/cauth"
	"github.com/tusharsoni/copper/cerror"
	"github.com/tusharsoni/copper/chttp"
	"github.com/tusharsoni/copper/clogger"
	"github.com/tusharsoni/copper/cmailer"
	"go.uber.org/fx"
	"gorm.io/gorm"
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

func NewModule(p NewModuleParams) (*Module, error) {
	svc, err := NewSvc(p.Auth, NewSQLRepo(p.DB), p.Mailer, p.Config)
	if err != nil {
		return nil, cerror.New(err, "failed to create service", nil)
	}

	router := NewRouter(RouterParams{
		RW:     chttp.NewJSONReaderWriter(p.Logger),
		Logger: p.Logger,
		Auth:   svc,
	})

	return &Module{
		Svc:    svc,
		Routes: router.Routes(),
	}, nil
}

var Fx = fx.Provide(
	NewSQLRepo,
	NewSvc,

	NewRouter,
	NewLogin,
	NewSignup,
)

var FxMigrations = fx.Invoke(
	RunMigrations,
)
