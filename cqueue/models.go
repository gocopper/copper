package cqueue

import (
	"time"
)

const (
	TaskStatusQueued     = "QUEUED"
	TaskStatusProcessing = "PROCESSING"
	TaskStatusCompleted  = "COMPLETED"
	TaskStatusFailed     = "FAILED"
)

type Task struct {
	UUID      string    `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`

	Type    string `gorm:"not null"`
	Payload []byte
	Status  string `gorm:"not null"`
	Error   *string
	Result  []byte
}

func (t Task) TableName() string {
	return "cqueue"
}
