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
	if len(expired) == 0 {
		return nil
	}
	// 快照被 reap 的原状态，便于持久化失败时回滚。
	type reaped struct {
		id string
		j  job.Job
	}
	rolled := make([]reaped, 0, len(expired))
	for _, id := range expired {
		j := q.inflight[id]
		rolled = append(rolled, reaped{id: id, j: j})
		delete(q.inflight, id)
		j.Worker = ""
		j.State = job.Pending
		q.pending = append(q.pending, j)
	}
	if err := q.persistLocked(); err != nil {
		// 持久化失败：把被 reap 的任务还原回 inflight，回滚 pending。
		q.pending = q.pending[:len(q.pending)-len(rolled)]
		for k := len(rolled) - 1; k >= 0; k-- {
			r := rolled[k]
			q.inflight[r.id] = r.j
		}
		return err
	}
	for range rolled {
		q.metrics.IncReaped()
	}
	return nil
}
