package cauth

import (
	"gorm.io/gorm"
)

// NewMigration instantiates and creates a new migration for cauth models
func NewMigration(db *gorm.DB) *Migration {
	return &Migration{db: db}
}

// Migration can create db tables for cauth models
type Migration struct {
	db *gorm.DB
}

// Run creates the tables corresponding to cauth models using the given db connection.
func (m *Migration) Run() error {
	return m.db.AutoMigrate(User{}, Session{})
}
