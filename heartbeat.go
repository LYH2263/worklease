package worklease

import (
	"context"
	"time"

	"example.com/worklease/internal/job"
	"example.com/worklease/internal/lease"
)

func (q *Queue) Heartbeat(ctx context.Context, id, worker string) error {

	if err := lease.Wait(ctx, time.Millisecond); err != nil {
		return wrapCancel(err)
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return ErrClosed
	}
	j, ok := q.inflight[id]
	if !ok {
		return ErrNotFound
	}
	if j.Worker != worker {
		return ErrNotOwner
	}
	if q.clk.Now().After(j.LeaseUntil) {
		return ErrExpired
	}
	j.LeaseUntil = q.clk.Now().Add(q.opts.LeaseTTL)
	j.State = job.Inflight
	q.inflight[id] = j
	q.metrics.IncHeartbeat()
	return q.persistLocked()
}
