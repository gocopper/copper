package email

import "time"

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

func (Credentials) TableName() string {
	return "cauth_email_credentials"
}
