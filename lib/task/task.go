// Package task provides the units of work processed by the workers.
package task

import (
	"time"

	"ridesharing/lib/result"
)

// Task is any unit of work that a worker can process.
type Task interface {
	Process() *result.Result
	GetID() int
	GetStatus() TaskStatus
}

// Base holds the fields shared by every task; embed it in concrete tasks.
type Base struct {
	taskID    int
	status    TaskStatus
	createdAt time.Time
}

func NewBase(taskID int) Base {
	return Base{taskID: taskID, status: Pending, createdAt: time.Now()}
}

func (b *Base) GetID() int {
	return b.taskID
}

func (b *Base) GetStatus() TaskStatus {
	return b.status
}

func (b *Base) GetCreatedAt() time.Time {
	return b.createdAt
}
