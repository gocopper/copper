package cacl

import (
	"github.com/jinzhu/gorm"
	"github.com/tusharsoni/copper/cerror"
)

func runMigrations(db *gorm.DB) error {
	err := db.AutoMigrate(&permission{}, &role{}, &roleUserJoin{}).Error
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

	err = db.
		Model(&roleUserJoin{}).
		AddUniqueIndex("idx_cacl_role_user_joins_user_uuid_role_uuid", "user_uuid", "role_uuid").
		Error
	if err != nil {
		return cerror.New(err, "failed to add unique index to cacl_role_user_joins", map[string]string{
			"idx": "idx_cacl_role_user_joins_user_uuid_role_uuid",
		})
	}

	err = db.Model(&roleUserJoin{}).AddForeignKey("user_uuid", "cauth_users(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		return cerror.New(err, "failed to add foreign key to role_user_joins.user_uuid", nil)
	}

	err = db.Model(&roleUserJoin{}).AddForeignKey("role_uuid", "cacl_roles(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		return cerror.New(err, "failed to add foreign key to role_user_joins.role_uuid", nil)
	}

	return nil
}
