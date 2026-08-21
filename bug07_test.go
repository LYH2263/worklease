package worklease_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"example.com/worklease"
	"example.com/worklease/internal/clock"
)

func TestBug07_ClaimContextHonorsCancel(t *testing.T) {
	q, err := worklease.New(worklease.Options{Clock: clock.NewFake(time.Now())})
	if err != nil {
		t.Fatal(err)
	}
	defer q.Close()
	_ = q.Enqueue(worklease.EnqueueSpec{ID: "7", Payload: []byte("p")})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = q.ClaimContext(ctx, "w")
	if err == nil || (!errors.Is(err, worklease.ErrCanceled) && !errors.Is(err, context.Canceled)) {
		t.Fatalf("got %v", err)
	}
}
