package ctags

import "time"

type tag struct {
	ID uint `gorm:"primary_key"`

	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`

	Tag      string `gorm:"not null"`
	EntityID string `gorm:"not null"`
}

func (tag) TableName() string {
	return "ctags"
}
