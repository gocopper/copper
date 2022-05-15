package csql

import (
	"github.com/gocopper/copper/clogger"
)

// Migration can be implemented by any struct that runs a database migration
type Migration interface {
	Run() error
}

// NewMigratorParams holds the params needed for NewMigrator
type NewMigratorParams struct {
	Migrations []Migration
	Logger     clogger.Logger
}

// NewMigrator creates a new Migrator
func NewMigrator(p NewMigratorParams) *Migrator {
	return &Migrator{
		migrations: p.Migrations,
		logger:     p.Logger,
	}
}

// Migrator can run database migrations by running all of the provided migrations
type Migrator struct {
	migrations []Migration
	logger     clogger.Logger
}

// Run runs all of the provided database migrations
func (m *Migrator) Run() error {
	m.logger.Info("Running database migrations..")

	for _, cm := range m.migrations {
		err := cm.Run()
		if err != nil {
			return err
		}
	}

	return nil
}
