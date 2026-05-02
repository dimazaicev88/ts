package status

type Status string

const (
	Active    Status = "active"
	Pending   Status = "pending"
	Scheduled Status = "scheduled"
	Retry     Status = "retry"
	Archived  Status = "archived"
	Completed Status = "completed"
)

func (t Status) ToString() string {
	return string(t)
}
