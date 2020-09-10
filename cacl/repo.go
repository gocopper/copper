package cacl

import (
	"context"

	"github.com/tusharsoni/copper/cerror"
	"github.com/tusharsoni/copper/csql"
	"gorm.io/gorm"
)

type repo interface {
	AddPermission(ctx context.Context, p *permission) error
	DeletePermission(ctx context.Context, uuid string) error
	GetPermissionForGrantee(ctx context.Context, granteeID, resource, action string) (*permission, error)

	AddRole(ctx context.Context, r *role) error
	AddUserRole(ctx context.Context, userUUID, roleUUID string) error
	FindRolesForUserUUID(ctx context.Context, userUUID string) ([]role, error)
}

type sqlRepo struct {
	db *gorm.DB
}

func newSQLRepo(db *gorm.DB) repo {
	return &sqlRepo{
		db: db,
	}
}

func (r *sqlRepo) AddPermission(ctx context.Context, p *permission) error {
	err := csql.GetConn(ctx, r.db).Save(p).Error
	if err != nil {
		return cerror.New(err, "failed to upsert permission", nil)
	}

	return nil
}

func (r *sqlRepo) DeletePermission(ctx context.Context, uuid string) error {
	err := csql.GetConn(ctx, r.db).
		Delete(permission{}, permission{UUID: uuid}).
		Error
	if err != nil {
		return cerror.New(err, "failed to delete permission", map[string]interface{}{
			"uuid": uuid,
		})
	}

	return nil
}

func (r *sqlRepo) GetPermissionForGrantee(ctx context.Context, granteeID, resource, action string) (*permission, error) {
	var p permission

	err := csql.GetConn(ctx, r.db).
		Where(permission{GranteeID: granteeID, Resource: resource, Action: action}).
		Find(&p).
		Error
	if err != nil {
		return nil, cerror.New(err, "failed to query permission by grantee, resource, action", map[string]interface{}{
			"granteeID": granteeID,
			"resource":  resource,
			"action":    action,
		})
	}

	return &p, nil
}

func (r *sqlRepo) AddRole(ctx context.Context, rl *role) error {
	err := csql.GetConn(ctx, r.db).Save(rl).Error
	if err != nil {
		return cerror.New(err, "failed to upsert role", nil)
	}

	return nil
}

func (r *sqlRepo) AddUserRole(ctx context.Context, userUUID, roleUUID string) error {
	j := roleUserJoin{
		UserUUID: userUUID,
		RoleUUID: roleUUID,
	}

	err := csql.GetConn(ctx, r.db).Save(&j).Error
	if err != nil {
		return cerror.New(err, "failed to add user role join", map[string]interface{}{
			"userUUID": userUUID,
			"roleUUID": roleUUID,
		})
	}

	return nil
}

func (r *sqlRepo) FindRolesForUserUUID(ctx context.Context, userUUID string) ([]role, error) {
	var roles []role

	err := csql.GetConn(ctx, r.db).
		Model(&role{}).
		Joins("JOIN cacl_role_user_joins ON role_user_joins.role_uuid=cacl_roles.uuid").
		Where("user_uuid = ?", userUUID).
		Find(&roles).
		Error
	if err != nil {
		return nil, cerror.New(err, "failed to query roles by user uuid", map[string]interface{}{
			"userUUID": userUUID,
		})
	}

	return roles, nil
}
