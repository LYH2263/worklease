package worklease_test

import (
	"testing"
	"time"

	"example.com/worklease"
	"example.com/worklease/internal/clock"
)

func TestBug02_ListJobsSliceAlias(t *testing.T) {
	q, err := worklease.New(worklease.Options{Clock: clock.NewFake(time.Now())})
	if err != nil {
		t.Fatal(err)
	}
	defer q.Close()
	if err := q.Enqueue(worklease.EnqueueSpec{ID: "2", Payload: []byte("p"), Tags: []string{"keep"}}); err != nil {
		t.Fatal(err)
	}
	list := q.List()
	list[0].Tags[0] = "x"
	list2 := q.List()
	if list2[0].Tags[0] != "keep" {
		t.Fatalf("tags aliased: %v", list2[0].Tags)
	}
}
