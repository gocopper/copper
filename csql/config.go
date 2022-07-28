package csql

import (
	"strings"

	"github.com/gocopper/copper/cconfig"
	"github.com/gocopper/copper/cerrors"
	migrate "github.com/rubenv/sql-migrate"
)

// MigrationsSource are valid options for csql.migrations.source configuration option.
// Use "dir" to load migrations from the local filesystem.
// Use "embed" to load migrations from the embedded directory in the binary.
const (
	MigrationsSourceDir   = "dir"
	MigrationsSourceEmbed = "embed"
)

// MigrationsDirection are valid options for csql.migrations.direction configuration option.
// Use "up" when running forward migrations and "down" when rolling back migrations.
const (
	MigrationsDirectionUp   = "up"
	MigrationsDirectionDown = "down"
)

// LoadConfig loads the csql config from the app config
func LoadConfig(appConfig cconfig.Loader) (Config, error) {
	config := Config{
		Migrations: ConfigMigrations{
			Direction: MigrationsDirectionUp,
			Source:    MigrationsSourceEmbed,
		},
	}

	err := appConfig.Load("csql", &config)
	if err != nil {
		return Config{}, cerrors.New(err, "failed to load sql config", nil)
	}

	return config, nil
}

type (
	// Config configures the csql module
	Config struct {
		Dialect            string           `toml:"dialect"`
		DSN                string           `toml:"dsn"`
		Migrations         ConfigMigrations `toml:"migrations"`
		MaxOpenConnections *int             `toml:"max_open_connections"`
	}

	// ConfigMigrations configures the migrations
	ConfigMigrations struct {
		Direction string `toml:"direction"`
		Source    string `toml:"source"`
	}
)

func (cm ConfigMigrations) sqlMigrateDirection() (migrate.MigrationDirection, error) {
	switch strings.ToLower(cm.Direction) {
	case MigrationsDirectionUp:
		return migrate.Up, nil
	case MigrationsDirectionDown:
		return migrate.Down, nil
	default:
		return 0, cerrors.New(nil, "invalid migration direction", map[string]interface{}{
			"direction": cm.Direction,
		})
	}
}
