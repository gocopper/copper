package cauth

import "time"

type User struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`

	UUID         string `gorm:"unique_index;not null"`
	SessionToken string `gorm:"not null"`
}

func (User) TableName() string {
	return "cauth_users"
}
