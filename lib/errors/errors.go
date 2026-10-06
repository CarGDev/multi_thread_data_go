// Package errors provides the custom error types used by the ride sharing system.
package errors

import (
	stderrors "errors"
	"fmt"
)

func IsProcessingError(err error) bool {
	var target *ProcessingError
	return stderrors.As(err, &target)
}

func IsQueueError(err error) bool {
	var target *QueueError
	return stderrors.As(err, &target)
}

func IsFileIOError(err error) bool {
	var target *FileIOError
	return stderrors.As(err, &target)
}

// ProcessingError is raised when a task fails while being processed.
type ProcessingError struct {
	message string
	taskID  int
}

func NewProcessingError(taskID int, message string) *ProcessingError {
	return &ProcessingError{message: message, taskID: taskID}
}

func (e *ProcessingError) Error() string {
	return fmt.Sprintf("processing error (task %d): %s", e.taskID, e.message)
}

func (e *ProcessingError) GetMessage() string {
	return e.message
}

func (e *ProcessingError) GetTaskID() int {
	return e.taskID
}

// QueueError is raised when a task queue operation fails.
type QueueError struct {
	message string
}

func NewQueueError(message string) *QueueError {
	return &QueueError{message: message}
}

func (e *QueueError) Error() string {
	return "queue error: " + e.message
}

func (e *QueueError) GetMessage() string {
	return e.message
}

// FileIOError is raised when reading or writing a file fails.
type FileIOError struct {
	message string
}

func NewFileIOError(message string) *FileIOError {
	return &FileIOError{message: message}
}

func (e *FileIOError) Error() string {
	return "file I/O error: " + e.message
}

func (e *FileIOError) GetMessage() string {
	return e.message
}
