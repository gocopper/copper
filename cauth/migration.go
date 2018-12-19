package cauth

import (
	"github.com/tusharsoni/copper/clogger"

	"github.com/jinzhu/gorm"
)

func runMigrations(db *gorm.DB, logger clogger.Logger) error {
	logger.Info("Running cauth migrations..", nil)
	return db.AutoMigrate(user{}).Error
}
