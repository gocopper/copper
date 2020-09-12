package cqueue

import (
	"github.com/tusharsoni/copper/cerror"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	err := db.AutoMigrate(&Task{})
	if err != nil {
		return cerror.New(err, "failed to auto migrate cqueue models", nil)
	}

	return nil
}
