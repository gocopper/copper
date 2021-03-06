package main

import (
	cauthemail "github.com/tusharsoni/copper/cauth/email"
	"github.com/tusharsoni/copper/csql"
	"go.uber.org/fx"
)

var ConfigFx = fx.Provide(NewConfig)

type Config struct {
	fx.Out

	SQL       csql.Config
	AuthEmail cauthemail.Config
}

func NewConfig() Config {
	return Config{
		SQL: csql.Config{
			Host: "localhost",
			Port: 5432,
			Name: "copper_auth",
			User: "postgres",
		},
		AuthEmail: cauthemail.GetDefaultConfig(),
	}
}
