package worklease_test

import (
	"errors"
	"testing"
	"time"

	"example.com/worklease"
	"example.com/worklease/internal/clock"
)

func TestBug05_AckErrorWrapsSentinel(t *testing.T) {
	q, err := worklease.New(worklease.Options{Clock: clock.NewFake(time.Now())})
	if err != nil {
		t.Fatal(err)
	}
	defer q.Close()
	err = q.Ack("missing", "w")
	if !errors.Is(err, worklease.ErrNotFound) {
		t.Fatalf("got %v", err)
	}
}
