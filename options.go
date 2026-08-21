package worklease

import (
	"time"

	"example.com/worklease/internal/clock"
	"example.com/worklease/internal/policy"
)

type Options struct {
	Clock       clock.Clock
	LeaseTTL    time.Duration
	MaxDepth    int
	PersistPath string
	ReaperEvery time.Duration
	Policy      policy.Policy
}

func (o *Options) normalize() {
	if o.LeaseTTL <= 0 {
		o.LeaseTTL = 2 * time.Second
	}
	if o.MaxDepth <= 0 {
		o.MaxDepth = 1024
	}
	if o.ReaperEvery <= 0 {
		o.ReaperEvery = 200 * time.Millisecond
	}
	if o.Policy.LeaseTTL == 0 {
		o.Policy = policy.Default()
	}
}
