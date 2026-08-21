package worklease

import (
	"context"
	"testing"
	"time"

	"example.com/worklease/internal/clock"
)

// TestListAndClaimTagIsolation guards against the regression where CloneTags was a
// no-op and List's pending branch handed out raw aliases to the queue's stored
// slice headers. Desensitizing (or mutating) a returned view must not dirty the
// queue's real tags or payload — otherwise subsequent List/Claim observes
// polluted data and business validation fails / tasks get killed.
func TestListAndClaimTagIsolation(t *testing.T) {
	q, err := New(Options{
		Clock:       clock.Real{},
		LeaseTTL:    5 * time.Second,
		MaxDepth:    16,
		ReaperEvery: time.Hour,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(func() { _ = q.Close() })

	payload := []byte("p-original")
	tags := []string{"secret-1", "secret-2"}
	if err := q.Enqueue(EnqueueSpec{ID: "p1", Payload: payload, Tags: tags}); err != nil {
		t.Fatalf("Enqueue pending: %v", err)
	}

	// 1) Desensitize the EnqueueSpec's own tags slice after enqueue: must not
	//    leak into the stored pending job.
	tags[0] = "TAMPERED"
	payload[0] = 'X'

	listed := q.List()
	if got := listed[0].Tags; got[0] != "secret-1" || got[1] != "secret-2" {
		t.Fatalf("pending tags leaked after enqueue: got %v", got)
	}
	if got := string(listed[0].Payload); got != "p-original" {
		t.Fatalf("pending payload leaked after enqueue: got %q", got)
	}

	// 2) Desensitize a List view's tags/payload: must not dirty the stored job
	//    seen by the next List.
	listed[0].Tags[0] = "POISONED"
	listed[0].Payload[0] = 'Y'
	again := q.List()
	if got := again[0].Tags; got[0] != "secret-1" {
		t.Fatalf("List view mutated queue tags: got %v", got)
	}
	if got := string(again[0].Payload); got != "p-original" {
		t.Fatalf("List view mutated queue payload: got %q", got)
	}

	// 3) Claim returns its own copies; desensitizing the claimed view must not
	//    affect the inflight job or a later List.
	claimed, err := q.ClaimContext(context.Background(), "w1")
	if err != nil {
		t.Fatalf("Claim: %v", err)
	}
	claimed.Tags[0] = "CLAIMED-TAMPER"
	claimed.Payload[0] = 'Z'

	inflight := q.List()
	var iv JobView
	for _, v := range inflight {
		if v.ID == "p1" {
			iv = v
		}
	}
	if iv.Tags[0] != "secret-1" {
		t.Fatalf("Claim view mutated inflight tags: got %v", iv.Tags)
	}
	if got := string(iv.Payload); got != "p-original" {
		t.Fatalf("Claim view mutated inflight payload: got %q", got)
	}
}
