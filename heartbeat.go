package worklease

import (
	"context"
	"time"

	"example.com/worklease/internal/job"
	"example.com/worklease/internal/lease"
)

func (q *Queue) Heartbeat(ctx context.Context, id, worker string) error {
	if err := ctx.Err(); err != nil {
		return wrapCancel(err)
	}
	if err := q.pol.WaitHeartbeat(ctx); err != nil {
		return wrapCancel(err)
	}
	// 可选短等待（测试可取消）
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
	prevLease := j.LeaseUntil
	j.LeaseUntil = q.clk.Now().Add(q.opts.LeaseTTL)
	j.State = job.Inflight
	q.inflight[id] = j
	if err := q.persistLocked(); err != nil {
		// 持久化失败：还原旧租约，避免内存与磁盘不一致。
		j.LeaseUntil = prevLease
		q.inflight[id] = j
		return err
	}
	q.metrics.IncHeartbeat()
	return nil
}
