package cauth

import (
	"encoding/json"
	"time"
)

type user struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`

	UUID     string  `gorm:"unique_index;not null"`
	Email    *string `gorm:"unique_index"`
	Password *string

	VerificationCode string `gorm:"not null"`
	Verified         bool   `gorm:"not null;default:false"`

	SessionToken *string
	LastLoginAt  *time.Time
}

func (user) TableName() string {
	return "cauth_users"
}

func (u *user) MarshalJSON() ([]byte, error) {
	var user struct {
		Email    *string `json:"email,omitempty"`
		Verified bool    `json:"verified"`
	}

	user.Email = u.Email
	user.Verified = u.Verified

	return json.Marshal(user)
}
