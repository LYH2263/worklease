package worklease

import "example.com/worklease/internal/job"

func (q *Queue) Ack(id, worker string) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return ErrClosed
	}
	j, ok := q.inflight[id]
	if !ok {
		return wrapNotFound(id)
	}
	if j.Worker != worker {
		return ErrNotOwner
	}
	delete(q.inflight, id)
	if err := q.persistLocked(); err != nil {
		// 持久化失败：放回 inflight，保持内存与磁盘一致。
		q.inflight[id] = j
		return err
	}
	q.metrics.IncAcked()
	return nil
}

func (q *Queue) Nack(id, worker string) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return ErrClosed
	}
	j, ok := q.inflight[id]
	if !ok {
		return wrapNotFound(id)
	}
	if j.Worker != worker {
		return ErrNotOwner
	}
	prevLease := j.LeaseUntil
	delete(q.inflight, id)
	j.Worker = ""
	j.LeaseUntil = q.clk.Now().Add(0)
	j.State = job.Pending
	q.pending = append(q.pending, j)
	if err := q.persistLocked(); err != nil {
		// 持久化失败：还原 inflight，回滚 pending。
		q.pending = q.pending[:len(q.pending)-1]
		j.Worker = worker
		j.LeaseUntil = prevLease
		j.State = job.Inflight
		q.inflight[id] = j
		return err
	}
	q.metrics.IncNacked()
	return nil
}
