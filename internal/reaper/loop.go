package reaper

import "time"

type Loop struct {
	Every time.Duration
}

func New(every time.Duration) *Loop {
	return &Loop{Every: every}
}
