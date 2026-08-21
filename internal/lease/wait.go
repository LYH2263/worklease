package lease

import (
	"context"
	"time"
)

// Wait blocks for at most d, returning nil once the duration elapses.
// It returns ctx.Err() (context.Canceled or context.DeadlineExceeded)
// as soon as ctx is canceled, without waiting for d.
//
// Wait is the throttle gate used by Heartbeat in the renewal loop: it must
// honor both the requested spacing (d) and cancellation so that a canceled
// renewal returns promptly instead of holding the worker past its lease
// deadline, which would drain the renewal pool and let leases expire.
func Wait(ctx context.Context, d time.Duration) error {
	// Fast path: nothing to wait for, only honor cancellation.
	if d <= 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return nil
		}
	}

	// Already canceled/deadline-exceeded: return without allocating a timer.
	if err := ctx.Err(); err != nil {
		return err
	}

	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
