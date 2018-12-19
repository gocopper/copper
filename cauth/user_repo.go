package cauth

import (
	"context"

	"github.com/tusharsoni/copper/cerror"
	"github.com/tusharsoni/copper/csql"

	"github.com/jinzhu/gorm"
)

// ErrUserNotFound is returned when the user is not found when queried by a unique attribute such as id or email.
var ErrUserNotFound = gorm.ErrRecordNotFound

// UserRepo provides methods to query and update users.
type UserRepo interface {
	FindByEmail(ctx context.Context, email string) (*user, error)
	Add(ctx context.Context, user *user) error
}

type sqlUserRepo struct {
	db *gorm.DB
}

func newSQLUserRepo(db *gorm.DB) UserRepo {
	return &sqlUserRepo{
		db: db,
	}
}

func (r *sqlUserRepo) FindByEmail(ctx context.Context, email string) (*user, error) {
	var u user

	err := csql.GetConn(ctx, r.db).
		Where(user{Email: email}).
		Find(&u).
		Error
	if err != nil {
		return nil, cerror.New(err, "failed to query u by email", map[string]string{
			"email": email,
		})
	}

	return &u, nil
}

func (r *sqlUserRepo) Add(ctx context.Context, user *user) error {
	err := csql.GetConn(ctx, r.db).Save(user).Error
	if err != nil {
		return cerror.New(err, "failed to upsert user", nil)
	}

	return nil
}
