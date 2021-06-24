package cauth

import (
	"time"
)

// User represents a user who has created an account. This model stores their login credentials
// as well as their metadata.
type User struct {
	UUID      string    `gorm:"primaryKey" json:"uuid"`
	CreatedAt time.Time `gorm:"not null" json:"-"`
	UpdatedAt time.Time `gorm:"not null" json:"-"`

	Email    *string `json:"email,omitempty"`
	Username *string `json:"username,omitempty"`

	Password           []byte `json:"-"`
	PasswordResetToken []byte `json:"-"`
}

// TableName returns the table name where the users are stored.
func (u User) TableName() string {
	return "cauth_users"
}

// Session represents a single logged-in session that a user is able create after providing valid
// login credentials.
type Session struct {
	UUID      string    `gorm:"primaryKey" json:"uuid"`
	CreatedAt time.Time `gorm:"not null" json:"-"`

	UserUUID  string    `gorm:"not null" json:"user_uuid"`
	Token     []byte    `gorm:"not null" json:"-"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
}

// TableName returns the table name where the sessions are stored.
func (u Session) TableName() string {
	return "cauth_sessions"
}
