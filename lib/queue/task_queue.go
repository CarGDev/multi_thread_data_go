// Package queue provides the shared, channel-backed queue of tasks.
package queue

import (
	"sync"

	"ridesharing/lib/errors"
	"ridesharing/lib/task"
)

// TaskQueue is a bounded queue built on a buffered channel. The channel hands
// each task to exactly one worker; the RWMutex only guards the closed flag so
// Enqueue never sends on a closed channel.
type TaskQueue struct {
	tasks  chan task.Task
	lock   sync.RWMutex
	closed bool
}

func NewTaskQueue(capacity int) *TaskQueue {
	return &TaskQueue{tasks: make(chan task.Task, capacity)}
}

// Enqueue blocks while the queue is full. Workers must be running (or the
// queue large enough) before submitting more tasks than the capacity.
func (q *TaskQueue) Enqueue(t task.Task) error {
	q.lock.RLock()
	defer q.lock.RUnlock()

	if q.closed {
		return errors.NewQueueError("queue is closed")
	}
	q.tasks <- t
	return nil
}

// Dequeue blocks until a task is available. It returns a QueueError once the
// queue is closed and drained.
func (q *TaskQueue) Dequeue() (task.Task, error) {
	t, ok := <-q.tasks
	if !ok {
		return nil, errors.NewQueueError("queue is closed and empty")
	}
	return t, nil
}

// TryDequeue returns immediately; ok is false when no task is available.
func (q *TaskQueue) TryDequeue() (t task.Task, ok bool) {
	select {
	case t, ok = <-q.tasks:
		return t, ok
	default:
		return nil, false
	}
}

func (q *TaskQueue) Size() int {
	return len(q.tasks)
}

func (q *TaskQueue) IsEmpty() bool {
	return len(q.tasks) == 0
}

// Close stops accepting tasks; queued tasks can still be dequeued.
func (q *TaskQueue) Close() {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.closed {
		return
	}
	q.closed = true
	close(q.tasks)
}

func (q *TaskQueue) IsClosed() bool {
	q.lock.RLock()
	defer q.lock.RUnlock()

	return q.closed
}
