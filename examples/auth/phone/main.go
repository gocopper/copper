package main

import (
	"github.com/tusharsoni/copper"
	"github.com/tusharsoni/copper/cauth"
	cauthphone "github.com/tusharsoni/copper/cauth/phone"
	"github.com/tusharsoni/copper/clogger"
	"github.com/tusharsoni/copper/csql"
	"github.com/tusharsoni/copper/ctexter"
)

func main() {
	app := copper.NewHTTPApp(
		clogger.StdFx,

		ConfigFx,

		csql.Fx,
		ctexter.FxLogger,

		cauth.Fx,
		cauth.FxMigrations,

		cauthphone.Fx,
		cauthphone.FxMigrations,
		cauthphone.FxValidators,
	)

	app.Run()
}
