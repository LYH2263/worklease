package worklease_test

import (
	"context"
	"testing"
	"time"

	"example.com/worklease"
	"example.com/worklease/internal/clock"
)

func TestBug01_ClaimPayloadSliceAlias(t *testing.T) {
	clk := clock.NewFake(time.Now())
	q, err := worklease.New(worklease.Options{Clock: clk})
	if err != nil {
		t.Fatal(err)
	}
	defer q.Close()
	payload := []byte("HELLO")
	if err := q.Enqueue(worklease.EnqueueSpec{ID: "1", Payload: payload, Tags: []string{"a"}}); err != nil {
		t.Fatal(err)
	}
	view, err := q.ClaimContext(context.Background(), "w1")
	if err != nil {
		t.Fatal(err)
	}
	view.Payload[0] = 'X'
	list := q.List()
	var found bool
	for _, j := range list {
		if j.ID == "1" {
			found = true
			if string(j.Payload) != "HELLO" {
				t.Fatalf("payload aliased: %q", j.Payload)
			}
		}
	}
	if !found {
		t.Fatal("missing")
	}
}
