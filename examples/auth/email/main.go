package main

import (
	"github.com/tusharsoni/copper"
	"github.com/tusharsoni/copper/cauth"
	cauthemail "github.com/tusharsoni/copper/cauth/email"
	"github.com/tusharsoni/copper/clogger"
	"github.com/tusharsoni/copper/cmailer"
	"github.com/tusharsoni/copper/csql"
)

func main() {
	app := copper.NewHTTPApp(
		clogger.StdFx,

		ConfigFx,

		csql.Fx,
		cmailer.LoggerFX,

		cauth.Fx,
		cauth.FxMigrations,

		cauthemail.Fx,
		cauthemail.FxMigrations,
	)

	app.Run()
}
