package cacl

import "time"

type permission struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	UUID      string `gorm:"uniqueIndex"`
	GranteeID string `gorm:"uniqueIndex:idx_grantee_id_resource_action"` // can be user or role uuid
	Resource  string `gorm:"uniqueIndex:idx_grantee_id_resource_action"`
	Action    string `gorm:"uniqueIndex:idx_grantee_id_resource_action"`
}

func (permission) TableName() string {
	return "cacl_permissions"
}

type role struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	UUID string `gorm:"uniqueIndex"`
	Name string
}

func (role) TableName() string {
	return "cacl_roles"
}

type roleUserJoin struct {
	ID uint `gorm:"primaryKey"`

	CreatedAt time.Time

	UserUUID string `gorm:"uniqueIndex:idx_user_uuid_role_uuid"`
	RoleUUID string `gorm:"uniqueIndex:idx_user_uuid_role_uuid"`
}

func (roleUserJoin) TableName() string {
	return "cacl_role_user_joins"
}
