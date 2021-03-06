package email

import (
	"context"

	"github.com/tusharsoni/copper/cerror"
	"github.com/tusharsoni/copper/csql"
	"gorm.io/gorm"
)

type Repo interface {
	GetCredentialsByUserUUID(ctx context.Context, userUUID string) (*Credentials, error)
	GetCredentialsByEmail(ctx context.Context, email string) (*Credentials, error)
	AddCredentials(ctx context.Context, c *Credentials) error
}

type sqlRepo struct {
	db *gorm.DB
}

func NewSQLRepo(db *gorm.DB) Repo {
	return &sqlRepo{
		db: db,
	}
}

func (r *sqlRepo) GetCredentialsByUserUUID(ctx context.Context, userUUID string) (*Credentials, error) {
	var c Credentials

	err := csql.GetConn(ctx, r.db).
		Where(Credentials{UserUUID: userUUID}).
		First(&c).
		Error
	if err != nil {
		return nil, cerror.New(err, "failed to query credentials", map[string]interface{}{
			"userUUID": userUUID,
		})
	}

	return &c, nil
}

func (r *sqlRepo) GetCredentialsByEmail(ctx context.Context, email string) (*Credentials, error) {
	var c Credentials

	err := csql.GetConn(ctx, r.db).
		Where(Credentials{Email: email}).
		First(&c).
		Error
	if err != nil {
		return nil, cerror.New(err, "failed to query credentials", map[string]interface{}{
			"email": email,
		})
	}

	return &c, nil
}

func (r *sqlRepo) AddCredentials(ctx context.Context, c *Credentials) error {
	err := csql.GetConn(ctx, r.db).Save(c).Error
	if err != nil {
		return cerror.New(err, "failed to upsert credentials", nil)
	}

	return nil
}
