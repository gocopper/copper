package cauth

import (
	"gorm.io/gorm"
)

// NewMigrator instantiates and creates a new migrator for cauth models
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{db: db}
}

// Migrator can create db tables for cauth models
type Migrator struct {
	db *gorm.DB
}

// Run creates the tables corresponding to cauth models using the given db connection.
func (m *Migrator) Run() error {
	return m.db.AutoMigrate(User{}, Session{})
}
