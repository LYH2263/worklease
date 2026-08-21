package worklease_test

import (
	"context"
	"testing"
	"time"

	"example.com/worklease"
	"example.com/worklease/internal/clock"
)

func TestBug08_HeartbeatWaitHonorsContext(t *testing.T) {
	clk := clock.NewFake(time.Now())
	q, err := worklease.New(worklease.Options{Clock: clk, LeaseTTL: time.Minute})
	if err != nil {
		t.Fatal(err)
	}
	defer q.Close()
	_ = q.Enqueue(worklease.EnqueueSpec{ID: "8", Payload: []byte("p")})
	_, err = q.ClaimContext(context.Background(), "w")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	start := time.Now()
	err = q.Heartbeat(ctx, "8", "w")
	if err == nil {
		t.Fatal("expected cancel")
	}
	if time.Since(start) > 2*time.Second {
		t.Fatalf("heartbeat ignored ctx: %v", time.Since(start))
	}
}
