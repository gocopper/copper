package cacl

import "time"

type permission struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time

	UUID      string `gorm:"unique_index"`
	GranteeID string // can be user or role uuid
	Resource  string
	Action    string
}

func (permission) TableName() string {
	return "cacl_permissions"
}

type role struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time

	UUID string `gorm:"unique_index"`
	Name string
}

func (role) TableName() string {
	return "cacl_roles"
}

type roleUserJoin struct {
	ID uint `gorm:"primary_key"`

	CreatedAt time.Time

	UserUUID string
	RoleUUID string
}

func (roleUserJoin) TableName() string {
	return "cacl_role_user_joins"
}
