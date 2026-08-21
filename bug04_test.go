package worklease_test

import (
	"context"
	"errors"
	"testing"

	"example.com/worklease"
)

func TestBug04_NilClockNoPanic(t *testing.T) {
	q, err := worklease.New(worklease.Options{Clock: nil})
	if err != nil {
		t.Fatal(err)
	}
	defer q.Close()
	_ = q.Enqueue(worklease.EnqueueSpec{ID: "4", Payload: []byte("p")})
	_, err = q.ClaimContext(context.Background(), "w")
	if !errors.Is(err, worklease.ErrNoClock) {
		t.Fatalf("got %v", err)
	}
}
