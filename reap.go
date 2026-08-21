package worklease

import "example.com/worklease/internal/job"

func (q *Queue) ReapExpired() error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed || q.clk == nil {
		return nil
	}
	now := q.clk.Now()
	var expired []string
	for id, j := range q.inflight {
		if now.After(j.LeaseUntil) {
			expired = append(expired, id)
		}
	}
	for _, id := range expired {
		j := q.inflight[id]
		delete(q.inflight, id)
		j.Worker = ""
		j.State = job.Pending
		q.pending = append(q.pending, j)
		q.metrics.IncReaped()
	}
	if len(expired) > 0 {
		return q.persistLocked()
	}
	return nil
}
