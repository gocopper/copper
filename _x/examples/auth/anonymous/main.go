package main

import (
	"github.com/tusharsoni/copper"
	"github.com/tusharsoni/copper/cauth"
	cauthanonymous "github.com/tusharsoni/copper/cauth/anonymous"
	"github.com/tusharsoni/copper/clogger"
	"github.com/tusharsoni/copper/cmailer"
	"github.com/tusharsoni/copper/csql"
	"go.uber.org/fx"
)

func main() {
	app := copper.NewHTTPApp(
		clogger.StdFx,

		ConfigFx,

		csql.Fx,
		cmailer.LoggerFX,

		cauth.Fx,
		cauth.FxMigrations,

		cauthanonymous.Fx,

		fx.Provide(
			NewRouter,
			NewProtectedRoute,
		),
	)

	app.Run()
}
