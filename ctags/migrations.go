package ctags

import (
	"github.com/jinzhu/gorm"
	"github.com/tusharsoni/copper/cerror"
)

func runMigrations(db *gorm.DB) error {
	err := db.AutoMigrate(&tag{}).Error
	if err != nil {
		return cerror.New(err, "failed to auto migrate tags models", nil)
	}

	err = db.
		Model(&tag{}).
		AddUniqueIndex("idx_ctags_tag_entity_id", "tag", "entity_id").
		Error
	if err != nil {
		return cerror.New(err, "failed to add unique index to tags", map[string]interface{}{
			"idx": "idx_ctags_tag_entity_id",
		})
	}

	return nil
}
