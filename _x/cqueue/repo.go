package cqueue

import (
	"context"

	"github.com/tusharsoni/copper/cerror"
	"github.com/tusharsoni/copper/csql"
	"gorm.io/gorm"
)

type Repo interface {
	AddTask(ctx context.Context, task *Task) error
	GetAndUpdateTaskForProcessing(ctx context.Context, taskType string) (*Task, error)
	GetTaskByUUID(ctx context.Context, uuid string) (*Task, error)
}

func NewSQLRepo(db *gorm.DB) Repo {
	return &sqlRepo{
		db: db,
	}
}

type sqlRepo struct {
	db *gorm.DB
}

func (r *sqlRepo) AddTask(ctx context.Context, task *Task) error {
	err := csql.GetConn(ctx, r.db).Save(task).Error
	if err != nil {
		return cerror.New(err, "failed to add task", nil)
	}

	return nil
}

func (r *sqlRepo) GetAndUpdateTaskForProcessing(ctx context.Context, taskType string) (*Task, error) {
	var (
		task  Task
		query = `
				UPDATE cqueue SET status=?
				WHERE uuid = (
					SELECT uuid FROM cqueue
					WHERE status=?
					AND type=?
					ORDER BY created_at ASC
					FOR UPDATE SKIP LOCKED 
					LIMIT 1
				)
				RETURNING *;`
	)

	err := csql.GetConn(ctx, r.db).
		Raw(query, TaskStatusProcessing, TaskStatusQueued, taskType).
		Scan(&task).
		Error
	if err != nil {
		return nil, cerror.New(err, "failed to query and update task", map[string]interface{}{
			"taskType": taskType,
		})
	}

	if task.UUID == "" {
		return nil, nil
	}

	return &task, nil
}

func (r *sqlRepo) GetTaskByUUID(ctx context.Context, uuid string) (*Task, error) {
	var task Task

	err := csql.GetConn(ctx, r.db).
		Where(&Task{UUID: uuid}).
		Find(&task).
		Error
	if err != nil {
		return nil, cerror.New(err, "failed to query task", map[string]interface{}{
			"uuid": uuid,
		})
	}

	return &task, nil
}
