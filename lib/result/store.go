package result

import (
	"fmt"
	"os"
	"sync"
	"time"

	"ridesharing/lib/errors"
)

type ResultStore struct {
	results []*Result
	lock    sync.Mutex
}

func NewResultStore() *ResultStore {
	return &ResultStore{}
}

func (s *ResultStore) AddResult(r *Result) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.results = append(s.results, r)
}

// GetResults returns a copy of the stored results.
func (s *ResultStore) GetResults() []*Result {
	s.lock.Lock()
	defer s.lock.Unlock()

	return append([]*Result(nil), s.results...)
}

func (s *ResultStore) Count() int {
	s.lock.Lock()
	defer s.lock.Unlock()

	return len(s.results)
}

func (s *ResultStore) CountSuccess() int {
	s.lock.Lock()
	defer s.lock.Unlock()

	n := 0
	for _, r := range s.results {
		if r.success {
			n++
		}
	}
	return n
}

func (s *ResultStore) CountFailed() int {
	return s.Count() - s.CountSuccess()
}

func (s *ResultStore) WriteToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return errors.NewFileIOError(err.Error())
	}
	defer file.Close()

	for _, r := range s.GetResults() {
		_, err := fmt.Fprintf(file, "task=%d success=%t completedAt=%s message=%s\n",
			r.taskID, r.success, r.completedAt.Format(time.RFC3339), r.message)
		if err != nil {
			return errors.NewFileIOError(err.Error())
		}
	}
	return nil
}
