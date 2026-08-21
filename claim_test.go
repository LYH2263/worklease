package worklease

import (
	"context"
	"errors"
	"testing"
	"time"

	"example.com/worklease/internal/clock"
)

// TestClaimContextCanceledDoesNotDispatch 锁定修复：传入已取消的 ctx 时，
// ClaimContext 不得出队、不得占住租约、不得把任务交给 worker。
// 修复前 WaitClaim 丢弃 ctx，已取消的领取仍会拿到任务并进入 inflight，
// 超时重试会把同一 payload 派给两个 worker。
func TestClaimContextCanceledDoesNotDispatch(t *testing.T) {
	q, err := New(Options{Clock: clock.NewFake(time.Unix(0, 0))})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := q.Enqueue(EnqueueSpec{ID: "j1", Payload: []byte("p1")}); err != nil {
		t.Fatalf("Enqueue: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	jv, err := q.ClaimContext(ctx, "w1")
	if err == nil {
		t.Fatalf("ClaimContext on canceled ctx should error, got job %q", jv.ID)
	}
	if !errors.Is(err, ErrCanceled) {
		t.Fatalf("expected ErrCanceled, got %v", err)
	}

	// 任务必须仍在 pending，未被占住租约、未分配 worker。
	if len(q.inflight) != 0 {
		t.Fatalf("canceled claim must not create inflight lease, got %d", len(q.inflight))
	}
	if len(q.pending) != 1 || q.pending[0].ID != "j1" {
		t.Fatalf("canceled claim must leave job pending, got pending=%v", q.pending)
	}
	if q.pending[0].Worker != "" {
		t.Fatalf("canceled claim must not assign worker, got %q", q.pending[0].Worker)
	}
	if q.pending[0].Attempts != 0 {
		t.Fatalf("canceled claim must not bump attempts, got %d", q.pending[0].Attempts)
	}
}

// TestClaimContextLiveStillWorks 保证修复没有误伤正常路径：
// 未取消的 ctx 仍能正常领取任务并占住租约。
func TestClaimContextLiveStillWorks(t *testing.T) {
	q, err := New(Options{Clock: clock.NewFake(time.Unix(0, 0))})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := q.Enqueue(EnqueueSpec{ID: "j1", Payload: []byte("p1")}); err != nil {
		t.Fatalf("Enqueue: %v", err)
	}

	jv, err := q.ClaimContext(context.Background(), "w1")
	if err != nil {
		t.Fatalf("ClaimContext: %v", err)
	}
	if jv.ID != "j1" || jv.Worker != "w1" {
		t.Fatalf("unexpected job view: %+v", jv)
	}
	if len(q.inflight) != 1 {
		t.Fatalf("live claim should create inflight lease, got %d", len(q.inflight))
	}
	if len(q.pending) != 0 {
		t.Fatalf("live claim should drain pending, got %d", len(q.pending))
	}
}
