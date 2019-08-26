package ctags

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/tusharsoni/copper/cerror"
	"github.com/tusharsoni/copper/csql"
)

type repo interface {
	Add(ctx context.Context, t *tag) error
	FindByEntityID(ctx context.Context, entityID string) ([]tag, error)
	FindByTag(ctx context.Context, tag string) ([]tag, error)
}

type sqlRepo struct {
	db *gorm.DB
}

func newSQLRepo(db *gorm.DB) repo {
	return &sqlRepo{
		db: db,
	}
}

func (r *sqlRepo) Add(ctx context.Context, t *tag) error {
	err := csql.GetConn(ctx, r.db).FirstOrCreate(t, &tag{
		Tag:      t.Tag,
		EntityID: t.EntityID,
	}).Error
	if err != nil {
		return cerror.New(err, "failed to upsert tag", nil)
	}

	return nil
}

func (r *sqlRepo) FindByEntityID(ctx context.Context, entityID string) ([]tag, error) {
	var tags []tag

	err := csql.GetConn(ctx, r.db).
		Model(&tag{}).
		Where(tag{EntityID: entityID}).
		Find(&tags).
		Error
	if err != nil {
		return nil, cerror.New(err, "failed to query tags by entity id", map[string]interface{}{
			"entityID": entityID,
		})
	}

	return tags, nil
}

func (r *sqlRepo) FindByTag(ctx context.Context, t string) ([]tag, error) {
	var tags []tag

	err := csql.GetConn(ctx, r.db).
		Model(&tag{}).
		Where(tag{Tag: t}).
		Find(&tags).
		Error
	if err != nil {
		return nil, cerror.New(err, "failed to query tags by tag", map[string]interface{}{
			"tag": t,
		})
	}

	return tags, nil
}
