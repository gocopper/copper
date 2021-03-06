package csql

import (
	"context"

	"github.com/tusharsoni/copper"
	"github.com/tusharsoni/copper/cconfig"
	"github.com/tusharsoni/copper/cerrors"
	"github.com/tusharsoni/copper/clogger"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// NewDBConnection creates and returns a new database connection. The connection is closed when the app exits.
func NewDBConnection(lc *copper.Lifecycle, appConfig cconfig.Config, logger clogger.Logger) (*gorm.DB, error) {
	var config struct {
		DSN string
	}

	err := appConfig.Load("csql", &config)
	if err != nil {
		return nil, cerrors.New(err, "failed to load sql config", nil)
	}

	dialect := newDialect(config.DSN)

	if config.DSN == "" {
		return nil, cerrors.New(err, "csql.dsn is not set in the config", map[string]interface{}{
			"dialect": dialect.Name(),
		})
	}

	logger.Info("Opening a database connection..")

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
