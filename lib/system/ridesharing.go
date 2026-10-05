// Package system wires the queue, workers, results and logger together.
package system

import (
	"sync"

	"ridesharing/lib/logger"
	"ridesharing/lib/queue"
	"ridesharing/lib/result"
	"ridesharing/lib/task"
	"ridesharing/lib/worker"
)

type RideSharingSystem struct {
	taskQueue   *queue.TaskQueue
	resultStore *result.ResultStore
	logger      *logger.Logger
	workers     []*worker.Worker
	wg          sync.WaitGroup
	running     bool
}

func Initialize(workerCount int) *RideSharingSystem {
	s := &RideSharingSystem{
		taskQueue:   queue.NewTaskQueue(),
		resultStore: result.NewResultStore(),
		logger:      logger.NewLogger(),
	}
	for i := 1; i <= workerCount; i++ {
		s.workers = append(s.workers, worker.NewWorker(i, s.taskQueue, s.resultStore, s.logger))
	}
	return s
}

func (s *RideSharingSystem) SubmitTask(t task.Task) error {
	if err := s.taskQueue.Enqueue(t); err != nil {
		s.logger.Error(err.Error())
		return err
	}
	return nil
}

func (s *RideSharingSystem) StartWorkers() {
	if s.running {
		return
	}
	s.running = true
	for _, w := range s.workers {
		w.Start(&s.wg)
	}
}

// Shutdown closes the queue; workers finish the remaining tasks and exit.
func (s *RideSharingSystem) Shutdown() {
	s.taskQueue.Close()
	s.running = false
}

func (s *RideSharingSystem) AwaitCompletion() {
	s.wg.Wait()
}

func (s *RideSharingSystem) GetResults() []*result.Result {
	return s.resultStore.GetResults()
}

func (s *RideSharingSystem) WriteResults(path string) error {
	err := s.resultStore.WriteToFile(path)
	if err != nil {
		s.logger.Error(err.Error())
	}
	return err
}
