package chttp

type Config struct {
	Port       uint
	HealthPath string
}

func GetDefaultConfig() Config {
	return Config{
		Port:       7450,
		HealthPath: "/_health",
	}
}
