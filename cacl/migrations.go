package cacl

import (
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(&permission{}, &role{}, &roleUserJoin{})
}
