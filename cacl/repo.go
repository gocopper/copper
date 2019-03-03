package cacl

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/tusharsoni/copper/cerror"
	"github.com/tusharsoni/copper/csql"
)

type repo interface {
	AddPermission(ctx context.Context, p *permission) error
	DeletePermission(ctx context.Context, uuid string) error
	GetPermissionForGrantee(ctx context.Context, granteeID, resource, action string) (*permission, error)
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
		return cerror.New(err, "failed to delete permission", map[string]string{
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
		return nil, cerror.New(err, "failed to query permission by grantee, resource, action", map[string]string{
			"granteeID": granteeID,
			"resource":  resource,
			"action":    action,
		})
	}

	return &p, nil
}
