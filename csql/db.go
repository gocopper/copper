package csql

import (
	"context"

	"github.com/gocopper/copper/cerrors"
	"github.com/gocopper/copper/clifecycle"
	"github.com/gocopper/copper/clogger"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// NewDBConnection creates and returns a new database connection. The connection is closed when the app exits.
func NewDBConnection(lc *clifecycle.Lifecycle, config Config, logger clogger.Logger) (*gorm.DB, error) {
	var dialect gorm.Dialector

	switch config.Dialect {
	case "sqlite":
		dialect = sqlite.Open(config.DSN)
	case "postgres":
		dialect = postgres.Open(config.DSN)
	default:
		return nil, cerrors.New(nil, "unknown dialect", map[string]interface{}{
			"dialect": config.Dialect,
		})
	}

	logger.WithTags(map[string]interface{}{
		"dialect": config.Dialect,
	}).Info("Opening a database connection..")

	db, err := gorm.Open(dialect, &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		return nil, cerrors.New(err, "failed to open db connection", nil)
	}

	lc.OnStop(func(ctx context.Context) error {
		logger.Info("Closing database connection..")

		sqlDB, err := db.DB()
		if err != nil {
			return cerrors.New(err, "failed to get *sql.DB from gorm db connection", nil)
		}

		return sqlDB.Close()
	})

	return db, nil
}
