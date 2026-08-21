package job

import "time"

type State string

const (
	Pending  State = "pending"
	Inflight State = "inflight"
	Done     State = "done"
)

type Job struct {
	ID         string
	Payload    []byte
	Tags       []string
	Worker     string
	LeaseUntil time.Time
	Attempts   int
	State      State
}
