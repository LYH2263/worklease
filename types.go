package worklease

import "time"

type JobView struct {
	ID       string
	Payload  []byte
	Tags     []string
	Worker   string
	LeaseUntil time.Time
	Attempts int
	State    string
}

type EnqueueSpec struct {
	ID      string
	Payload []byte
	Tags    []string
}
