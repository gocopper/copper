package cqueue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/tusharsoni/copper/cerror"
	"go.uber.org/fx"
)

type QueueTaskParams struct {
	Type    string
	Payload json.RawMessage
}

type Svc interface {
	QueueTask(ctx context.Context, p QueueTaskParams) (*Task, error)
	DequeueTask(ctx context.Context, taskType string) (*Task, error)
	GetTask(ctx context.Context, uuid string) (*Task, error)
	FailTask(ctx context.Context, uuid, error string) error
	CompleteTask(ctx context.Context, uuid string, result []byte) error
}

type SvcParams struct {
	fx.In

	Repo   Repo
	Config Config
}

func NewSvc(p SvcParams) Svc {
	return &svc{
		repo:   p.Repo,
		config: p.Config,
	}
}

type svc struct {
	repo   Repo
	config Config
}

func (s *svc) QueueTask(ctx context.Context, p QueueTaskParams) (*Task, error) {
	task := Task{
		UUID:    uuid.New().String(),
		Type:    p.Type,
		Payload: p.Payload,
		Status:  TaskStatusQueued,
	}

	err := s.repo.AddTask(ctx, &task)
	if err != nil {
		return nil, cerror.New(err, "failed to save task", map[string]interface{}{
			"type": p.Type,
		})
	}

	return &task, nil
}

func (s *svc) DequeueTask(ctx context.Context, taskType string) (*Task, error) {
	for {
		task, err := s.repo.GetAndUpdateTaskForProcessing(ctx, taskType)
		if err != nil {
			return nil, err
		}

		if task == nil {
			time.Sleep(s.config.DequeueWaitTime)
			continue
		}

		return task, nil
	}
}

func (s *svc) GetTask(ctx context.Context, uuid string) (*Task, error) {
	return s.repo.GetTaskByUUID(ctx, uuid)
}

func (s *svc) FailTask(ctx context.Context, uuid, taskError string) error {
	task, err := s.repo.GetTaskByUUID(ctx, uuid)
	if err != nil {
		return cerror.New(err, "failed to get task", map[string]interface{}{
			"uuid": uuid,
		})
	}

	if task.Status != TaskStatusProcessing {
		return cerror.New(nil, "task is not processing", map[string]interface{}{
			"uuid": uuid,
		})
	}

	task.Status = TaskStatusFailed
	task.Error = &taskError

	err = s.repo.AddTask(ctx, task)
	if err != nil {
		return cerror.New(err, "failed to update task", map[string]interface{}{
			"uuid":      uuid,
			"taskError": taskError,
		})
	}

	return nil
}

func (s *svc) CompleteTask(ctx context.Context, uuid string, result []byte) error {
	task, err := s.repo.GetTaskByUUID(ctx, uuid)
	if err != nil {
		return cerror.New(err, "failed to get task", map[string]interface{}{
			"uuid": uuid,
		})
	}

	if task.Status != TaskStatusProcessing {
		return cerror.New(nil, "task is not processing", map[string]interface{}{
			"uuid": uuid,
		})
	}

	task.Status = TaskStatusCompleted
	task.Result = result

	err = s.repo.AddTask(ctx, task)
	if err != nil {
		return cerror.New(err, "failed to update task", map[string]interface{}{
			"uuid":   uuid,
			"result": string(result),
		})
	}

	return nil
}
