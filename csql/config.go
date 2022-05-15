package csql

import (
	"github.com/gocopper/copper/cconfig"
	"github.com/gocopper/copper/cerrors"
)

// LoadConfig loads the csql config from the app config
func LoadConfig(appConfig cconfig.Loader) (Config, error) {
	var config Config

	err := appConfig.Load("csql", &config)
	if err != nil {
		return Config{}, cerrors.New(err, "failed to load sql config", nil)
	}

	return config, nil
}

// Config configures the csql module
type Config struct {
	Dialect string `toml:"dialect"`
	DSN     string `toml:"dsn"`
}
