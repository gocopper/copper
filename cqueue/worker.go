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

func (w *Worker) run(ctx context.Context, s Svc, logger clogger.Logger) {
	var (
		limiter      <-chan time.Time
		workerLogger = logger.WithTags(map[string]interface{}{
			"taskType": w.TaskType,
		})
	)

	if w.RateLimit == nil {
		limiter = time.Tick(time.Nanosecond)
	} else {
		limiter = time.Tick(*w.RateLimit)
	}

	workerLogger.Info("Starting worker..")
	defer workerLogger.Info("Exiting worker..")

	for {
		select {
		case <-limiter:
			task, err := s.DequeueTask(ctx, w.TaskType)
			if err != nil {
				workerLogger.Error("Failed to dequeue task", err)
				return
			}

			taskLogger := workerLogger.WithTags(map[string]interface{}{
				"taskUUID": task.UUID,
			})
			ctx, cancel := context.WithTimeout(ctx, w.Timeout)

			err = w.Handler(ctx, task.Payload)
			if err != nil {
				taskLogger.Error("Failed to run task", err)

				err = s.UpdateTaskStatus(ctx, task.UUID, TaskStatusFailed, cptr.String(err.Error()))
				if err != nil {
					taskLogger.Error("Failed to mark task as failed", err)
				}

				cancel()
				continue
			}

			err = s.UpdateTaskStatus(ctx, task.UUID, TaskStatusCompleted, nil)
			if err != nil {
				taskLogger.Error("Failed to mark task as completed", err)
			}

			cancel()
		}
	}
}
