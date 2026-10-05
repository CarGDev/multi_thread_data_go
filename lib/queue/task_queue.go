// Package queue provides a blocking, thread-safe queue of tasks.
package queue

import (
	"sync"

	"ridesharing/lib/errors"
	"ridesharing/lib/task"
)

type TaskQueue struct {
	tasks         []task.Task
	lock          sync.Mutex
	taskAvailable *sync.Cond
	closed        bool
}

func NewTaskQueue() *TaskQueue {
	q := &TaskQueue{}
	q.taskAvailable = sync.NewCond(&q.lock)
	return q
}

func (q *TaskQueue) Enqueue(t task.Task) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.closed {
		return errors.NewQueueError("queue is closed")
	}
	q.tasks = append(q.tasks, t)
	q.taskAvailable.Signal()
	return nil
}

// Dequeue blocks until a task is available. It returns a QueueError once the
// queue is closed and drained.
func (q *TaskQueue) Dequeue() (task.Task, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	for len(q.tasks) == 0 && !q.closed {
		q.taskAvailable.Wait()
	}
	if len(q.tasks) == 0 {
		return nil, errors.NewQueueError("queue is closed and empty")
	}

	t := q.tasks[0]
	q.tasks = q.tasks[1:]
	return t, nil
}

func (q *TaskQueue) IsEmpty() bool {
	q.lock.Lock()
	defer q.lock.Unlock()

	return len(q.tasks) == 0
}

func (q *TaskQueue) Close() {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.closed = true
	q.taskAvailable.Broadcast()
}

func (q *TaskQueue) IsClosed() bool {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.closed
}
