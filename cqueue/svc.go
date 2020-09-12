package cqueue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/tusharsoni/copper/cerror"
	"github.com/tusharsoni/copper/clogger"
	"go.uber.org/fx"
)

type QueueTaskParams struct {
	Type    string
	Payload json.RawMessage
}

type Svc interface {
	QueueTask(ctx context.Context, p QueueTaskParams) error
	DequeueTask(ctx context.Context, taskType string) (*Task, error)
	UpdateTaskStatus(ctx context.Context, uuid, status string, taskErr *string) error
	StartWorkers(ctx context.Context)
}

type SvcParams struct {
	fx.In

	Repo    Repo
	Workers []Worker `group:"cqueue/workers"`
	Config  Config
	Logger  clogger.Logger
}

func NewSvc(p SvcParams) Svc {
	return &svc{
		repo:    p.Repo,
		workers: p.Workers,
		config:  p.Config,
		logger:  p.Logger,
	}
}

type svc struct {
	repo    Repo
	workers []Worker
	config  Config
	logger  clogger.Logger
}

func (s *svc) QueueTask(ctx context.Context, p QueueTaskParams) error {
	task := Task{
		UUID:    uuid.New().String(),
		Type:    p.Type,
		Payload: p.Payload,
		Status:  TaskStatusQueued,
	}

	return s.repo.AddTask(ctx, &task)
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

func (s *svc) UpdateTaskStatus(ctx context.Context, uuid, status string, taskErr *string) error {
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

	if status != TaskStatusCompleted && status != TaskStatusFailed {
		return cerror.New(nil, "invalid status", map[string]interface{}{
			"uuid":   uuid,
			"status": status,
		})
	}

	task.Status = status
	task.Error = taskErr

	err = s.repo.AddTask(ctx, task)
	if err != nil {
		return cerror.New(err, "failed to update task", map[string]interface{}{
			"uuid":      uuid,
			"status":    status,
			"taskError": taskErr,
		})
	}

	return nil
}

func (s *svc) StartWorkers(ctx context.Context) {
	for i := range s.workers {
		go s.workers[i].run(ctx, s, s.logger)
	}
}
