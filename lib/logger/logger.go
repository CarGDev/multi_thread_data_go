// Package logger provides a thread-safe logger.
package logger

import (
	"fmt"
	"os"
	"sync"
	"time"

	"ridesharing/lib/errors"
)

type Logger struct {
	lock sync.Mutex
	file *os.File
}

func NewLogger() *Logger {
	return &Logger{}
}

// NewFileLogger logs to stdout and also appends every line to the file at path.
func NewFileLogger(path string) (*Logger, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, errors.NewFileIOError(err.Error())
	}
	return &Logger{file: file}, nil
}

func (l *Logger) Close() error {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.file == nil {
		return nil
	}
	if err := l.file.Close(); err != nil {
		return errors.NewFileIOError(err.Error())
	}
	l.file = nil
	return nil
}

func (l *Logger) log(level, message string) {
	l.lock.Lock()
	defer l.lock.Unlock()

	line := fmt.Sprintf("%s [%s] %s\n", time.Now().Format(time.RFC3339), level, message)
	fmt.Print(line)
	if l.file != nil {
		l.file.WriteString(line)
	}
}

func (l *Logger) Info(message string) {
	l.log("INFO", message)
}

func (l *Logger) Warn(message string) {
	l.log("WARN", message)
}

func (l *Logger) Error(message string) {
	l.log("ERROR", message)
}

func (l *Logger) LogWorkerStart(workerID int) {
	l.Info(fmt.Sprintf("worker %d started", workerID))
}

func (l *Logger) LogWorkerComplete(workerID int) {
	l.Info(fmt.Sprintf("worker %d completed", workerID))
}

func (l *Logger) LogTaskStart(workerID, taskID int) {
	l.Info(fmt.Sprintf("worker %d started task %d", workerID, taskID))
}

func (l *Logger) LogTaskComplete(workerID, taskID int) {
	l.Info(fmt.Sprintf("worker %d completed task %d", workerID, taskID))
}

func (l *Logger) LogTaskError(workerID, taskID int, message string) {
	l.Error(fmt.Sprintf("worker %d task %d failed: %s", workerID, taskID, message))
}

func (l *Logger) LogException(workerID int, err error) {
	l.Error(fmt.Sprintf("worker %d: %v", workerID, err))
}
