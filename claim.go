package worklease

import (
	"context"

	"example.com/worklease/internal/job"
)

func (q *Queue) Claim(worker string, lease interface{}) (JobView, error) {
	return q.ClaimContext(context.Background(), worker)
}

func (q *Queue) ClaimContext(ctx context.Context, worker string) (JobView, error) {

	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return JobView{}, ErrClosed
	}
	if q.clk == nil {
		return JobView{}, ErrNoClock
	}
	if worker == "" {
		return JobView{}, ErrInvalid
	}
	// 已取消的 ctx 不得继续发放任务：否则取消信号被忽略，
	// 超时重试会把同一 payload 派给两个 worker。
	if err := ctx.Err(); err != nil {
		return JobView{}, wrapCancel(err)
	}
	if err := q.pol.WaitClaim(ctx); err != nil {
		return JobView{}, wrapCancel(err)
	}
	// WaitClaim 可能阻塞；唤醒后再次确认 ctx 仍未被取消，
	// 避免在已取消路径上继续出队、占住租约。
	if err := ctx.Err(); err != nil {
		return JobView{}, wrapCancel(err)
	}
	if len(q.pending) == 0 {
		return JobView{}, ErrEmpty
	}
	j := q.pending[0]
	q.pending = q.pending[1:]
	j.State = job.Inflight
	j.Worker = worker
	j.Attempts++
	j.LeaseUntil = q.clk.Now().Add(q.opts.LeaseTTL)
	q.inflight[j.ID] = j
	if err := q.persistLocked(); err != nil {
		delete(q.inflight, j.ID)
		q.pending = append([]job.Job{j}, q.pending...)
		return JobView{}, err
	}
	q.metrics.IncClaimed()
	return JobView{
		ID:         j.ID,
		Payload:    job.CloneBytes(j.Payload),
		Tags:       job.CloneTags(j.Tags),
		Worker:     j.Worker,
		LeaseUntil: j.LeaseUntil,
		Attempts:   j.Attempts,
		State:      string(j.State),
	}, nil
}
