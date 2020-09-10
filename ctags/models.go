package ctags

import "time"

type tag struct {
	ID uint `gorm:"primaryKey"`

	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`

	Tag      string `gorm:"not null;uniqueIndex:idx_tag_entity_id"`
	EntityID string `gorm:"not null;uniqueIndex:idx_tag_entity_id"`
}

func (tag) TableName() string {
	return "ctags"
}
