package csql

// Config holds the parameters to connect to a Postgres database. Provide this config using Fx.
type Config struct {
	Host     string
	Port     uint
	Name     string
	User     string
	Password string
}
