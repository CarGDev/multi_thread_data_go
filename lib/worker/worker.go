// Package worker provides the goroutine workers that process queued tasks.
package worker

import (
	"fmt"
	"sync"
	"sync/atomic"

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
	defer func() {
		if r := recover(); r != nil {
			err := errors.NewProcessingError(t.GetID(), fmt.Sprint(r))
			w.logger.LogException(w.workerID, err)
			w.resultStore.AddResult(result.NewResult(t.GetID(), false, err.GetMessage()))
		}
	}()

	w.resultStore.AddResult(t.Process())
}

func (w *Worker) Stop() {
	w.running.Store(false)
}
