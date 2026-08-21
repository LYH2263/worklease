package worklease

import "example.com/worklease/internal/job"

func (q *Queue) Enqueue(spec EnqueueSpec) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return ErrClosed
	}
	if spec.ID == "" || spec.Payload == nil {
		return ErrInvalid
	}
	if len(q.pending)+len(q.inflight) >= q.opts.MaxDepth {
		return ErrInvalid
	}
	j := job.Job{
		ID:       spec.ID,
		Payload:  job.CloneBytes(spec.Payload),
		Tags:     job.CloneTags(spec.Tags),
		State:    job.Pending,
		Attempts: 0,
	}
	q.pending = append(q.pending, j)
	if err := q.persistLocked(); err != nil {
		// 持久化失败：回滚 pending，避免内存有、磁盘没有的「幽灵任务」。
		q.pending = q.pending[:len(q.pending)-1]
		return err
	}
	q.metrics.IncEnqueued()
	return nil
}
