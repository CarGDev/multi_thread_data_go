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

const (
	queueCapacity = 100
	logFilePath   = "ridesharing.log"
)

func Initialize(workerCount int) *RideSharingSystem {
	log, err := logger.NewFileLogger(logFilePath)
	if err != nil {
		log = logger.NewLogger()
		log.Warn("file logging disabled: " + err.Error())
	}

	s := &RideSharingSystem{
		taskQueue:   queue.NewTaskQueue(queueCapacity),
		resultStore: result.NewResultStore(),
		logger:      log,
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

// SubmitTasks enqueues every task and returns the first error, if any.
func (s *RideSharingSystem) SubmitTasks(tasks []task.Task) error {
	for _, t := range tasks {
		if err := s.SubmitTask(t); err != nil {
			return err
		}
	}
	return nil
}

type Stats struct {
	Total   int
	Success int
	Failed  int
}

func (s *RideSharingSystem) Stats() Stats {
	return Stats{
		Total:   s.resultStore.Count(),
		Success: s.resultStore.CountSuccess(),
		Failed:  s.resultStore.CountFailed(),
	}
}

// Run starts the workers, processes all tasks, waits for completion and
// writes the results to outPath.
func (s *RideSharingSystem) Run(tasks []task.Task, outPath string) error {
	s.StartWorkers()
	if err := s.SubmitTasks(tasks); err != nil {
		s.Shutdown()
		s.AwaitCompletion()
		return err
	}
	s.Shutdown()
	s.AwaitCompletion()
	return s.WriteResults(outPath)
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

// Close releases the log file; call it after AwaitCompletion.
func (s *RideSharingSystem) Close() error {
	return s.logger.Close()
}

func (s *RideSharingSystem) AwaitCompletion() {
	s.wg.Wait()
}

func (s *RideSharingSystem) GetResults() []*result.Result {
	return s.resultStore.GetResults()
}

func (s *RideSharingSystem) WriteResults(path string) error {
	err := s.resultStore.WriteToCSV(path)
	if err != nil {
		s.logger.Error(err.Error())
	}
	return err
}
