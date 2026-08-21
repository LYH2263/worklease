package worklease

import "example.com/worklease/internal/job"

func (q *Queue) List() []JobView {
	q.mu.Lock()
	defer q.mu.Unlock()
	out := make([]JobView, 0, len(q.pending)+len(q.inflight))
	for _, j := range q.pending {
		out = append(out, JobView{
			ID: j.ID, Payload: job.CloneBytes(j.Payload), Tags: job.CloneTags(j.Tags),
			State: string(j.State), Attempts: j.Attempts,
		})
	}
	for _, j := range q.inflight {
		out = append(out, JobView{
			ID: j.ID, Payload: job.CloneBytes(j.Payload), Tags: job.CloneTags(j.Tags),
			Worker: j.Worker, LeaseUntil: j.LeaseUntil, State: string(j.State), Attempts: j.Attempts,
		})
	}
	return out
}

func (q *Queue) Stats() map[string]int64 {
	q.mu.Lock()
	defer q.mu.Unlock()
	s := q.metrics.Snapshot()
	s["pending"] = int64(len(q.pending))
	s["inflight"] = int64(len(q.inflight))
	return s
}
