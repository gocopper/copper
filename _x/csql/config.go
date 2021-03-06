package csql

// Config holds the parameters to connect to a Postgres database
type Config struct {
	Host     string
	Port     uint
	Name     string
	User     string
	Password string
}
