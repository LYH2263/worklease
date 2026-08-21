package worklease_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"example.com/worklease"
	"example.com/worklease/internal/clock"
)

func TestBug03_ClaimAfterCloseNoPanic(t *testing.T) {
	q, err := worklease.New(worklease.Options{Clock: clock.NewFake(time.Now())})
	if err != nil {
		t.Fatal(err)
	}
	_ = q.Enqueue(worklease.EnqueueSpec{ID: "3", Payload: []byte("p")})
	_ = q.Close()
	_, err = q.ClaimContext(context.Background(), "w")
	if !errors.Is(err, worklease.ErrClosed) {
		t.Fatalf("got %v", err)
	}
}
