package worklease

import (
	"testing"
	"time"

	"example.com/worklease/internal/clock"
	"example.com/worklease/internal/policy"
)

// TestClaimPayloadIsolation guards against the Claim→job.CloneBytes aliasing
// bug: a worker mutating the returned Payload bytes must not bleed back into the
// queue's internal inflight table, which is what List() reads.
func TestClaimPayloadIsolation(t *testing.T) {
	q, err := New(Options{
		Clock:    clock.NewFake(time.Unix(0, 0)),
		Policy:   policy.Default(),
		LeaseTTL: time.Minute,
		MaxDepth: 16,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	original := []byte("hello world")
	if err := q.Enqueue(EnqueueSpec{ID: "j1", Payload: original}); err != nil {
		t.Fatalf("Enqueue: %v", err)
	}

	got, err := q.Claim("worker-A", nil)
	if err != nil {
		t.Fatalf("Claim: %v", err)
	}

	// Mutate the returned payload the way a worker stamps a debug mark.
	got.Payload[0] = 'X'

	// The caller's own view reflects their own mutation...
	if string(got.Payload) != "Xello world" {
		t.Fatalf("caller payload = %q, want %q", got.Payload, "Xello world")
	}

	// ...but the queue's inflight entry must be untouched.
	view := q.List()
	var inflight *JobView
	for i := range view {
		if view[i].ID == "j1" {
			inflight = &view[i]
		}
	}
	if inflight == nil {
		t.Fatalf("j1 missing from List()")
	}
	if string(inflight.Payload) != "hello world" {
		t.Fatalf("inflight payload aliased with Claim result: got %q, want %q",
			inflight.Payload, "hello world")
	}

	// And the originally enqueued slice must be untouched too.
	if string(original) != "hello world" {
		t.Fatalf("source payload mutated: got %q, want %q", original, "hello world")
	}
}
