package worklease_test

import (
	"path/filepath"
	"testing"
	"time"

	"example.com/worklease"
	"example.com/worklease/internal/clock"
)

func TestBug10_CloseStopsReaperBeforeFlush(t *testing.T) {
	path := filepath.Join(t.TempDir(), "q.json")
	clk := clock.NewFake(time.Now())
	q, err := worklease.New(worklease.Options{Clock: clk, PersistPath: path, ReaperEvery: time.Hour})
	if err != nil {
		t.Fatal(err)
	}
	q.StartReaper()
	if err := q.Enqueue(worklease.EnqueueSpec{ID: "keep", Payload: []byte("p"), Tags: []string{"t"}}); err != nil {
		t.Fatal(err)
	}
	if err := q.Close(); err != nil {
		t.Fatal(err)
	}
	q2, err := worklease.New(worklease.Options{Clock: clk, PersistPath: path})
	if err != nil {
		t.Fatal(err)
	}
	defer q2.Close()
	if err := q2.LoadPersist(); err != nil {
		t.Fatal(err)
	}
	if len(q2.List()) < 1 {
		t.Fatal("Close flushed empty queue")
	}
}
