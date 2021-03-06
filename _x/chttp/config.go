package chttp

// Config can be used to customize the setup for chttp server.
// The value of this type should be provided as an fx.Option when creating a copper app.
type Config struct {
	Port       uint
	HealthPath string
}

func (c Config) isValid() bool {
	return c.Port > 0
}

// GetDefaultConfig provides the default config that can be used with chttp.
func GetDefaultConfig() Config {
	return Config{
		Port:       7450,
		HealthPath: "/_health",
	}
}
