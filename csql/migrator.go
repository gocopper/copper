package csql

import (
	"crypto/sha256"
	"database/sql"
	"embed"
	"fmt"
	"io"

	"github.com/gocopper/copper/cerrors"
	"github.com/gocopper/copper/clogger"
	migrate "github.com/rubenv/sql-migrate"
)

// Migrations is a collection of .sql files that represent the database schema
type Migrations embed.FS

// NewMigratorParams holds the params needed for NewMigrator
type NewMigratorParams struct {
	DB         *sql.DB
	Migrations Migrations
	Config     Config
	Logger     clogger.Logger
}

// NewMigrator creates a new Migrator
func NewMigrator(p NewMigratorParams) *Migrator {
	return &Migrator{
		db:         p.DB,
		migrations: embed.FS(p.Migrations),
		config:     p.Config,
		logger:     p.Logger,
	}
}

// Migrator can run database migrations by running the provided migrations in the migrations dir
type Migrator struct {
	db         *sql.DB
	migrations embed.FS
	config     Config
	logger     clogger.Logger
}

// Run runs the provided database migrations
func (m *Migrator) Run() error {
	m.logger.WithTags(map[string]interface{}{
		"direction": m.config.Migrations.Direction,
		"source":    m.config.Migrations.Source,
	}).Info("Running database migrations..")

	direction, err := m.config.Migrations.sqlMigrateDirection()
	if err != nil {
		return cerrors.New(err, "failed to get sql migrate direction from config", nil)
	}

	hasMigrations, err := m.hasMigrations()
	if err != nil {
		return cerrors.New(err, "failed to check for migrations", nil)
	}

	if !hasMigrations {
		m.logger.Info("No migrations found")
		return nil
	}

	source := migrate.MigrationSource(migrate.EmbedFileSystemMigrationSource{
		FileSystem: m.migrations,
		Root:       ".",
	})
	if m.config.Migrations.Source == MigrationsSourceDir {
		source = migrate.FileMigrationSource{
			Dir: "./migrations",
		}
	}

	migrateMax := 0 // no limit
	if direction == migrate.Down {
		migrateMax = 1 // only run 1 migration when reverting
	}

	dialect := m.config.Dialect
	if dialect == "pgx" {
		dialect = "postgres"
	}

	n, err := migrate.ExecMax(m.db, dialect, source, direction, migrateMax)
	if err != nil {
		return cerrors.New(err, "failed to exec database migrations", nil)
	}

	m.logger.WithTags(map[string]interface{}{
		"count": n,
	}).Info("Successfully applied migrations")

	return nil
}

// hasMigrations returns true if the migrations directory has at least 1 non-empty migration file.
func (m *Migrator) hasMigrations() (bool, error) {
	const emptyMigrationsChecksum = "fba9ab24993a94e181dc952f2568a4e98b47e331d89772af3115fe1c7b90d27f"

	entries, err := m.migrations.ReadDir(".")
	if err != nil {
		return false, cerrors.New(err, "failed to read migrations dir", nil)
	}

	if len(entries) == 0 {
		return false, nil
	}

	if len(entries) > 1 {
		return true, nil
	}

	f, err := m.migrations.Open(entries[0].Name())
	if err != nil {
		return false, cerrors.New(err, "failed to open migrations file", map[string]interface{}{
			"name": entries[0].Name(),
		})
	}
	defer func() { _ = f.Close() }()

	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		return false, cerrors.New(err, "failed to calculate sha256 for migration file", nil)
	}

	checksum := fmt.Sprintf("%x", h.Sum(nil))
	if checksum == emptyMigrationsChecksum {
		return false, nil
	}

	return true, nil
}
