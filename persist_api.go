package worklease

import (
	"example.com/worklease/internal/job"
	"example.com/worklease/internal/persist"
)

func (q *Queue) persistLocked() error {
	if q.persist == nil {
		return nil
	}
	snap := persist.Snapshot{
		Pending:  append([]job.Job(nil), q.pending...),
		Inflight: map[string]job.Job{},
	}
	for k, v := range q.inflight {
		snap.Inflight[k] = v
	}
	if err := q.persist.Save(snap); err != nil {
		return wrapPersist(err)
	}
	// Close 须释放 Save 打开的句柄
	return nil
}

func (q *Queue) LoadPersist() error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.persist == nil {
		return nil
	}
	snap, err := q.persist.Load()
	if err != nil {
		return wrapPersist(err)
	}
	q.pending = append([]job.Job(nil), snap.Pending...)
	q.inflight = make(map[string]job.Job, len(snap.Inflight))
	for k, v := range snap.Inflight {
		q.inflight[k] = v
	}
	return nil
}
