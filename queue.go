package worklease

import (
	"sync"

	"example.com/worklease/internal/clock"
	"example.com/worklease/internal/job"
	"example.com/worklease/internal/metrics"
	"example.com/worklease/internal/persist"
	"example.com/worklease/internal/policy"
	"example.com/worklease/internal/reaper"
)

type Queue struct {
	mu       sync.Mutex
	opts     Options
	pending  []job.Job
	inflight map[string]job.Job
	clk      clock.Clock
	pol      policy.Policy
	persist  *persist.Store
	metrics  *metrics.Registry
	reaper   *reaper.Loop
	closed   bool
	stopCh   chan struct{}
	doneCh   chan struct{}
}

func New(opts Options) (*Queue, error) {
	opts.normalize()
	q := &Queue{
		opts:     opts,
		inflight: make(map[string]job.Job),
		clk:      opts.Clock,
		pol:      opts.Policy,
		metrics:  metrics.New(),
		stopCh:   make(chan struct{}),
		doneCh:   make(chan struct{}),
	}
	close(q.doneCh)
	if opts.PersistPath != "" {
		q.persist = persist.New(opts.PersistPath)
	}
	return q, nil
}
