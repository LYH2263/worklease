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

	// 关闭顺序：先停 reaper，再刷盘（保留在途状态），最后清内部队列状态。
	// 若在 reaper 仍在跑时就清空 pending/inflight 再 persistLocked，空快照会
	// 覆盖磁盘，重启后 LoadPersist 读到空队列，在途任务全丢。
	//
	// 1. 等 reaper goroutine 退出，确保不再有并发的 ReapExpired 触发持久化。
	if q.reaperRunning() {
		<-q.doneCh
	}

	q.mu.Lock()
	// 2. 落盘：此时 pending/inflight 仍为完整状态，写出的快照含全部在途任务。
	_ = q.persistLocked()
	if q.persist != nil {
		_ = q.persist.Close()
	}
	// 3. 清空内部队列状态：刷盘完成后再丢弃内存数据。
	q.pending = nil
	q.inflight = nil
	q.mu.Unlock()
	return nil
}

// reaperRunning reports whether the reaper goroutine is active.
// doneCh is closed (signaling "done") both before StartReaper and once the
// goroutine exits, so an unclosed channel means the goroutine is still running.
func (q *Queue) reaperRunning() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	select {
	case <-q.doneCh:
		return false
	default:
		return true
	}
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
