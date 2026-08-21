package worklease

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"example.com/worklease/internal/clock"
)

// TestClaimContextAfterClose exercises the rolling-shutdown race:
// once Close returns, lingering Claim calls must return ErrClosed
// instead of touching the now-nil clk/pending/inflight and crashing workd.
func TestClaimContextAfterClose(t *testing.T) {
	q, err := New(Options{
		Clock:    clock.NewFake(time.Now()),
		LeaseTTL: 5 * time.Second,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := q.Enqueue(EnqueueSpec{ID: "j1", Payload: []byte("p")}); err != nil {
		t.Fatalf("Enqueue: %v", err)
	}
	if err := q.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	// Sanity: closed state means the live fields were hollowed out.
	q.mu.Lock()
	hollowed := q.clk == nil && q.pending == nil && q.inflight == nil
	q.mu.Unlock()
	if !hollowed {
		t.Fatal("Close did not nil out clk/pending/inflight; test premise wrong")
	}

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := q.ClaimContext(context.Background(), "worker-1")
			if !errors.Is(err, ErrClosed) {
				t.Errorf("ClaimContext after Close: want ErrClosed, got %v", err)
			}
		}()
	}
	wg.Wait()
}

// TestClaimContextAfterCloseEmpty ensures the guard fires even on a queue
// that was closed with zero pending jobs (the q.pending[0] deref path).
func TestClaimContextAfterCloseEmpty(t *testing.T) {
	q, err := New(Options{
		Clock:    clock.NewFake(time.Now()),
		LeaseTTL: 5 * time.Second,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := q.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	_, err = q.ClaimContext(context.Background(), "worker-1")
	if !errors.Is(err, ErrClosed) {
		t.Fatalf("want ErrClosed, got %v", err)
	}
}
