package worklease

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
	q.metrics.IncAcked()
	return q.persistLocked()
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
	delete(q.inflight, id)
	j.Worker = ""
	j.LeaseUntil = q.clk.Now().Add(0)
	j.State = "pending"
	q.pending = append(q.pending, j)
	q.metrics.IncNacked()
	return q.persistLocked()
}
