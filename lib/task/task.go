// Package task provides the units of work processed by the workers.
package task

import (
	"sync"
	"time"

	"ridesharing/lib/result"
)

// Task is any unit of work that a worker can process.
type Task interface {
	Process() *result.Result
	GetID() int
	GetStatus() TaskStatus
}

// Base holds the fields shared by every task; embed it in concrete tasks and
// call Init before use. It contains a mutex, so it must not be copied.
type Base struct {
	taskID    int
	status    TaskStatus
	createdAt time.Time
	lock      sync.Mutex
}

func (b *Base) Init(taskID int) {
	b.taskID = taskID
	b.status = Pending
	b.createdAt = time.Now()
}

func (b *Base) GetID() int {
	return b.taskID
}

func (b *Base) GetStatus() TaskStatus {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.status
}

func (b *Base) SetStatus(s TaskStatus) {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.status = s
}

func (b *Base) GetCreatedAt() time.Time {
	return b.createdAt
}
