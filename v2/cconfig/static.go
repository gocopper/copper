package cconfig

// NewStaticConfig provides a config based on the given static config.
func NewStaticConfig(config map[string]interface{}) Config {
	return &static{config: config}
}

type static struct {
	config map[string]interface{}
}

func (c *static) Value(path string) interface{} {
	return c.config[path]
}
