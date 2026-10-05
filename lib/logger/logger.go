// Package logger provides a thread-safe logger.
package logger

import (
	"fmt"
	"sync"
	"time"
)

type Logger struct {
	lock sync.Mutex
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) log(level, message string) {
	l.lock.Lock()
	defer l.lock.Unlock()

	fmt.Printf("%s [%s] %s\n", time.Now().Format(time.RFC3339), level, message)
}

func (l *Logger) Info(message string) {
	l.log("INFO", message)
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

func (l *Logger) LogException(workerID int, err error) {
	l.Error(fmt.Sprintf("worker %d: %v", workerID, err))
}
