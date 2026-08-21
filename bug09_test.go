package worklease_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"example.com/worklease"
	"example.com/worklease/internal/clock"
)

func TestBug09_PersistFileClosed(t *testing.T) {
	path := filepath.Join(t.TempDir(), "q.json")
	q, err := worklease.New(worklease.Options{Clock: clock.NewFake(time.Now()), PersistPath: path})
	if err != nil {
		t.Fatal(err)
	}
	if err := q.Enqueue(worklease.EnqueueSpec{ID: "9", Payload: []byte("p")}); err != nil {
		t.Fatal(err)
	}
	if err := q.Close(); err != nil {
		t.Fatal(err)
	}
	rotated := path + ".1"
	if err := os.Rename(path, rotated); err != nil {
		t.Fatalf("rename after Close: %v", err)
	}
}
