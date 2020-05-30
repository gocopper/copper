package cauth2

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/tusharsoni/copper/cerror"
	"github.com/tusharsoni/copper/csql"
)

type Repo interface {
	GetUser(ctx context.Context, uuid string) (*User, error)
	AddUser(ctx context.Context, user *User) error
}

type sqlRepo struct {
	db *gorm.DB
}

func NewSQLRepo(db *gorm.DB) Repo {
	return &sqlRepo{
		db: db,
	}
}

func (r *sqlRepo) GetUser(ctx context.Context, uuid string) (*User, error) {
	var u User

	err := csql.GetConn(ctx, r.db).
		Where(User{UUID: uuid}).
		Find(&u).
		Error
	if err != nil {
		return nil, cerror.New(err, "failed to query user", map[string]interface{}{
			"uuid": uuid,
		})
	}

	return &u, nil
}

func (r *sqlRepo) AddUser(ctx context.Context, user *User) error {
	err := csql.GetConn(ctx, r.db).Save(user).Error
	if err != nil {
		return cerror.New(err, "failed to upsert user", nil)
	}

	return nil
}
