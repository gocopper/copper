package cqueue

import (
	"context"
	"time"

	"github.com/tusharsoni/copper/clogger"
	"github.com/tusharsoni/copper/cptr"
	"go.uber.org/fx"
)

type WorkerResult struct {
	fx.Out

	Worker Worker `group:"cqueue/workers"`
}

type Worker struct {
	TaskType  string
	Timeout   time.Duration
	RateLimit *time.Duration
	Handler   func(ctx context.Context, payload []byte) error
}

func (w *Worker) Start(ctx context.Context, s Svc, logger clogger.Logger) {
	var (
		limiter <-chan time.Time
		log     = logger.WithTags(map[string]interface{}{
			"taskType": w.TaskType,
		})
	)

	log.Info("Starting background worker..")
	defer log.Info("Exiting background worker..")

	for {
		if limiter == nil {
			w.runNextTask(ctx, s, log)
		} else {
			select {
			case <-limiter:
				w.runNextTask(ctx, s, log)
			}
		}
	}
}

func (w *Worker) runNextTask(ctx context.Context, s Svc, logger clogger.Logger) {
	task, err := s.DequeueTask(ctx, w.TaskType)
	if err != nil {
		logger.Error("Failed to dequeue task", err)
		return
	}

	log := logger.WithTags(map[string]interface{}{
		"taskUUID": task.UUID,
	})

	ctx, cancel := context.WithTimeout(ctx, w.Timeout)
	defer cancel()

	log.Info("Running task..")
	defer log.Info("Task completed")

	err = w.Handler(ctx, task.Payload)
	if err != nil {
		log.Error("Failed to run task", err)

		err = s.UpdateTaskStatus(ctx, task.UUID, TaskStatusFailed, cptr.String(err.Error()))
		if err != nil {
			log.Error("Failed to mark task as failed", err)
		}

		return
	}

	err = s.UpdateTaskStatus(ctx, task.UUID, TaskStatusCompleted, nil)
	if err != nil {
		log.Error("Failed to mark task as completed", err)
	}
}
