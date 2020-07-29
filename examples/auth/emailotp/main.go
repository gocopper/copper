package main

import (
	"github.com/tusharsoni/copper"
	"github.com/tusharsoni/copper/cauth"
	cauthemailotp "github.com/tusharsoni/copper/cauth/emailotp"
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

		cauthemailotp.Fx,
		cauthemailotp.FxMigrations,
	)

	app.Run()
}
