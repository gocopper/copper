package cauth

import (
	"github.com/tusharsoni/copper/v2/clogger"
	"gorm.io/gorm"
)

// Migrate creates the tables corresponding to cauth models using the given db connection.
func Migrate(db *gorm.DB, logger clogger.Logger) error {
	return db.AutoMigrate(User{}, Session{}) //nolint: exhaustivestruct
}
