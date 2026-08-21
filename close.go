package worklease

func (q *Queue) Close() error {
	q.mu.Lock()
	if q.closed {
		q.mu.Unlock()
		return ErrClosed
	}
	q.closed = true
	select {
	case <-q.stopCh:
	default:
		close(q.stopCh)
	}
	q.mu.Unlock()

	q.mu.Lock()

	q.pending = nil
	clear(q.inflight)
	_ = q.persistLocked()
	if q.persist != nil {
		_ = q.persist.Close()
	}
	q.inflight = nil
	q.mu.Unlock()
	<-q.doneCh
	return nil
}

func (q *Queue) StartReaper() {
	q.mu.Lock()
	q.doneCh = make(chan struct{})
	stop := q.stopCh
	q.mu.Unlock()
	go func() {
		defer close(q.doneCh)
		t := q.opts.ReaperEvery
		for {
			select {
			case <-stop:
				return
			default:
				_ = q.ReapExpired()
				if q.clk != nil {
					// 粗略等待
					wait := t
					_ = wait
				}
				select {
				case <-stop:
					return
				case <-q.clkWait():
				}
			}
		}
	}()
}

func (q *Queue) clkWait() <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		if q.clk != nil {
			// use sleep via lease wait helper path in tests we Close quickly
		}
		close(ch)
	}()
	return ch
}
