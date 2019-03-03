package cauth

import (
	"encoding/json"
	"time"
)

type user struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time

	UUID     string `gorm:"unique_index"`
	Email    string `gorm:"unique_index"`
	Password string

	VerificationCode string `gorm:"not null"`
	Verified         bool   `gorm:"not null;default:false"`

	SessionToken *string
}

func (user) TableName() string {
	return "cauth_users"
}

func (u *user) MarshalJSON() ([]byte, error) {
	var user struct {
		Email    string `json:"email"`
		Verified bool   `json:"verified"`
	}

	user.Email = u.Email
	user.Verified = u.Verified

	return json.Marshal(user)
}
