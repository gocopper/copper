package csql_test

import (
	"database/sql"
	"embed"
	"testing"

	"github.com/gocopper/copper/clogger"
	"github.com/gocopper/copper/csql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

//go:embed migrations_test.sql
var Migrations embed.FS

func TestMigrator_Run(t *testing.T) {
	t.Parallel()

	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	// migrate up

	migratorUp := csql.NewMigrator(csql.NewMigratorParams{
		DB:         db,
		Migrations: csql.Migrations(Migrations),
		Config: csql.Config{
			Dialect: "sqlite3",
			Migrations: csql.ConfigMigrations{
				Direction: csql.MigrationsDirectionUp,
				Source:    csql.MigrationsSourceEmbed,
			},
		},
		Logger: clogger.NewNoop(),
	})

	err = migratorUp.Run()
	assert.NoError(t, err)

	res, err := db.Query("select * from people")
	assert.NoError(t, err)
	assert.NoError(t, res.Err())

	assert.True(t, res.Next())

	// migrate down

	migratorDown := csql.NewMigrator(csql.NewMigratorParams{
		DB:         db,
		Migrations: csql.Migrations(Migrations),
		Config: csql.Config{
			Dialect: "sqlite3",
			Migrations: csql.ConfigMigrations{
				Direction: csql.MigrationsDirectionDown,
				Source:    csql.MigrationsSourceEmbed,
			},
		},
		Logger: clogger.NewNoop(),
	})

	err = migratorDown.Run()
	assert.NoError(t, err)

	_, err = db.Query("select * from people") //nolint:rowserrcheck
	assert.EqualError(t, err, "no such table: people")
}
