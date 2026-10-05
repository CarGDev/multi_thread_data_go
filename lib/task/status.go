package task

type TaskStatus int

const (
	Pending TaskStatus = iota
	Processing
	Completed
	Failed
)

var statusNames = map[TaskStatus]string{
	Pending:    "pending",
	Processing: "processing",
	Completed:  "completed",
	Failed:     "failed",
}

func (s TaskStatus) String() string {
	if status, ok := statusNames[s]; ok {
		return status
	}
	return "pending"
}
