package phone

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/tusharsoni/copper/cerror"
	"github.com/tusharsoni/copper/csql"
)

type Repo interface {
	GetCredentialsByPhoneNumber(ctx context.Context, phoneNumber string) (*Credentials, error)
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

func (r *sqlRepo) GetCredentialsByPhoneNumber(ctx context.Context, phoneNumber string) (*Credentials, error) {
	var c Credentials

	err := csql.GetConn(ctx, r.db).
		Where(Credentials{PhoneNumber: phoneNumber}).
		Find(&c).
		Error
	if err != nil {
		return nil, cerror.New(err, "failed to query credentials", map[string]interface{}{
			"phoneNumber": phoneNumber,
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
