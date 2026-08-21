package worklease

import (
	"context"

	"example.com/worklease/internal/job"
)

func (q *Queue) Claim(worker string, lease interface{}) (JobView, error) {
	return q.ClaimContext(context.Background(), worker)
}

func (q *Queue) ClaimContext(ctx context.Context, worker string) (JobView, error) {
	if err := ctx.Err(); err != nil {
		return JobView{}, wrapCancel(err)
	}
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.clk == nil {
		return JobView{}, ErrNoClock
	}
	if worker == "" {
		return JobView{}, ErrInvalid
	}
	if err := q.pol.WaitClaim(ctx); err != nil {
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
