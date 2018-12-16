package chttp

type Config struct {
	Port       uint
	HealthPath string
}

func (c Config) isValid() bool {
	return c.Port > 0
}

func GetDefaultConfig() Config {
	return Config{
		Port:       7450,
		HealthPath: "/_health",
	}
}
