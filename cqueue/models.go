package cqueue

import (
	"encoding/json"
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

func (t *Task) MarshalJSON() ([]byte, error) {
	var j struct {
		UUID      string  `json:"uuid"`
		CreatedAt int64   `json:"created_at"`
		UpdatedAt int64   `json:"updated_at"`
		Type      string  `json:"type"`
		Payload   string  `json:"payload"`
		Status    string  `json:"status"`
		Error     *string `json:"error"`
		Result    string  `json:"result"`
	}

	j.UUID = t.UUID
	j.CreatedAt = t.CreatedAt.Unix()
	j.UpdatedAt = t.CreatedAt.Unix()
	j.Type = t.Type
	j.Payload = string(t.Payload)
	j.Status = t.Status
	j.Error = t.Error
	j.Result = string(t.Result)

	return json.Marshal(j)
}

func (t Task) TableName() string {
	return "cqueue"
}
