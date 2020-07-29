package email

import (
	"encoding/json"
	"time"
)

type Credentials struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`

	UserUUID         string `gorm:"not null;unique_index"`
	Email            string `gorm:"not null;unique_index"`
	Password         string `gorm:"not null"`
	Verified         bool   `gorm:"not null;default:false"`
	VerificationCode string `gorm:"not null"`
}

func (c Credentials) MarshalJSON() ([]byte, error) {
	var j struct {
		UserUUID string `json:"user_uuid"`
		Email    string `json:"email"`
		Verified bool   `json:"verified"`
	}

	j.UserUUID = c.UserUUID
	j.Email = c.Email
	j.Verified = c.Verified

	return json.Marshal(j)
}

func (Credentials) TableName() string {
	return "cauth_email_credentials"
}
