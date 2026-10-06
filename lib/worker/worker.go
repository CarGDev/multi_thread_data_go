// Package worker provides the goroutine workers that process queued tasks.
package worker

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"ridesharing/lib/errors"
	"ridesharing/lib/logger"
	"ridesharing/lib/queue"
	"ridesharing/lib/result"
	"ridesharing/lib/task"
)

type Worker struct {
	workerID    int
	taskQueue   *queue.TaskQueue
	resultStore *result.ResultStore
	logger      *logger.Logger
	running     atomic.Bool
}

func NewWorker(
	workerID int,
	taskQueue *queue.TaskQueue,
	resultStore *result.ResultStore,
	logger *logger.Logger,
) *Worker {
	return &Worker{
		workerID:    workerID,
		taskQueue:   taskQueue,
		resultStore: resultStore,
		logger:      logger,
	}
}

// Start launches the worker in its own goroutine; wg is released when it exits.
func (w *Worker) Start(wg *sync.WaitGroup) {
	w.running.Store(true)
	wg.Add(1)
	go func() {
		defer wg.Done()
		w.Run()
	}()
}

func (w *Worker) Run() {
	w.logger.LogWorkerStart(w.workerID)
	defer w.logger.LogWorkerComplete(w.workerID)

	for w.running.Load() {
		t, err := w.taskQueue.Dequeue()
		if err != nil {
			return // queue closed and drained
		}
		w.ProcessTask(t)
	}
}

func (w *Worker) ProcessTask(t task.Task) {
	w.logger.LogTaskStart(w.workerID, t.GetID())
	startedAt := time.Now()

	defer func() {
		if r := recover(); r != nil {
			err := errors.NewProcessingError(t.GetID(), fmt.Sprint(r))
			w.logger.LogException(w.workerID, err)

			res := result.NewResult(t.GetID(), false, err.GetMessage())
			res.SetWorker(w.workerID)
			res.SetStartedAt(startedAt)
			w.resultStore.AddResult(res)
		}
	}()

	res := t.Process()
	res.SetWorker(w.workerID)
	res.SetStartedAt(startedAt)
	w.resultStore.AddResult(res)

	if res.IsSuccess() {
		w.logger.LogTaskComplete(w.workerID, t.GetID())
	} else {
		w.logger.LogTaskError(w.workerID, t.GetID(), res.GetMessage())
	}
}

func (w *Worker) Stop() {
	w.running.Store(false)
}

func (w *Worker) GetID() int {
	return w.workerID
}

func (w *Worker) IsRunning() bool {
	return w.running.Load()
}
