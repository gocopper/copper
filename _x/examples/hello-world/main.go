package main

import (
	"github.com/tusharsoni/copper"
	"github.com/tusharsoni/copper/clogger"
)

func main() {
	params := copper.HTTPAppParams{
		Logger: clogger.NewStdLogger(),
	}

	copper.RunHTTPApp(params)
}
