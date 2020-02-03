package main

import (
	"github.com/tusharsoni/copper"
	"github.com/tusharsoni/copper/clogger"
)

func main() {
	app := copper.NewHTTPApp(
		clogger.StdFx,
	)

	app.Run()
}
