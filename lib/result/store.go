package result

import (
	"fmt"
	"os"
	"strings"
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

func (s *ResultStore) WriteToFile(path string) error {
	var sb strings.Builder
	for _, r := range s.GetResults() {
		fmt.Fprintf(&sb, "task=%d success=%t completedAt=%s message=%s\n",
			r.taskID, r.success, r.completedAt.Format(time.RFC3339), r.message)
	}

	if err := os.WriteFile(path, []byte(sb.String()), 0o644); err != nil {
		return errors.NewFileIOError(err.Error())
	}
	return nil
}
