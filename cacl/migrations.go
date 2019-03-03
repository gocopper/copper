package cacl

import (
	"github.com/jinzhu/gorm"
	"github.com/tusharsoni/copper/cerror"
)

func runMigrations(db *gorm.DB) error {
	err := db.AutoMigrate(&permission{}, &role{}).Error
	if err != nil {
		return cerror.New(err, "failed to auto migrate cacl models", nil)
	}

	err = db.
		Model(&permission{}).
		AddUniqueIndex("idx_cacl_permissions_grantee_id_resource_action", "grantee_id", "resource", "action").
		Error
	if err != nil {
		return cerror.New(err, "failed to add unique index to cacl_permissions", map[string]string{
			"idx": "idx_cacl_permissions_grantee_id_resource_action",
		})
	}

	return nil
}
