// Package result provides the outcome of a processed task and a store for them.
package result

import "time"

type Result struct {
	taskID      int
	success     bool
	message     string
	completedAt time.Time
}

func NewResult(taskID int, success bool, message string) *Result {
	return &Result{
		taskID:      taskID,
		success:     success,
		message:     message,
		completedAt: time.Now(),
	}
}

func (r *Result) GetTaskID() int {
	return r.taskID
}

func (r *Result) IsSuccess() bool {
	return r.success
}

func (r *Result) GetMessage() string {
	return r.message
}

func (r *Result) GetCompletedAt() time.Time {
	return r.completedAt
}
