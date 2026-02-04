package csql

import (
	"context"
	"database/sql"
	"time"

	"github.com/gocopper/copper/cerrors"
	"github.com/gocopper/copper/clifecycle"
	"github.com/gocopper/copper/clogger"
)

// NewDBConnection creates and returns a new database connection. The connection is closed when the app exits.
func NewDBConnection(lc *clifecycle.Lifecycle, config Config, logger clogger.Logger) (*sql.DB, error) {
	const (
		DefaultMaxOpenConnections = 25
		DefaultMaxIdleConnections = 25
		DefaultConnMaxLifetime    = 5 * time.Minute
	)

	logger.WithTags(map[string]interface{}{
		"dialect": config.Dialect,
	}).Info("Opening a database connection..")

	db, err := sql.Open(config.Dialect, config.DSN)
	if err != nil {
		return nil, cerrors.New(err, "failed to open db connection", map[string]interface{}{
			"dialect": config.Dialect,
		})
	}

	if config.MaxOpenConnections == nil {
		db.SetMaxOpenConns(DefaultMaxOpenConnections)
	} else {
		db.SetMaxOpenConns(*config.MaxOpenConnections)
	}

	if config.MaxIdleConnections == nil {
		db.SetMaxIdleConns(DefaultMaxIdleConnections)
	} else {
		db.SetMaxIdleConns(*config.MaxIdleConnections)
	}

	if config.ConnMaxLifetimeMins == nil {
		db.SetConnMaxLifetime(DefaultConnMaxLifetime)
	} else {
		db.SetConnMaxLifetime(time.Duration(*config.ConnMaxLifetimeMins) * time.Minute)
	}

	if err := db.Ping(); err != nil {
		return nil, cerrors.New(err, "failed to ping db", nil)
	}

	lc.OnStop(func(ctx context.Context) error {
		logger.Info("Closing database connection..")

		err := db.Close()
		if err != nil {
			return cerrors.New(err, "failed to close db connection", nil)
		}

		return nil
	})

	return db, nil
}
