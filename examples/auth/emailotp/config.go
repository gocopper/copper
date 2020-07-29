package main

import (
	cauthemailotp "github.com/tusharsoni/copper/cauth/emailotp"
	"github.com/tusharsoni/copper/csql"
	"go.uber.org/fx"
)

var ConfigFx = fx.Provide(NewConfig)

type Config struct {
	fx.Out

	SQL          csql.Config
	AuthEmailOTP cauthemailotp.Config
}

func NewConfig() Config {
	return Config{
		SQL: csql.Config{
			Host: "localhost",
			Port: 5432,
			Name: "copper_auth",
			User: "postgres",
		},
		AuthEmailOTP: cauthemailotp.GetDefaultConfig(),
	}
}
