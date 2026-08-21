package worklease_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"example.com/worklease"
	"example.com/worklease/internal/clock"
)

func TestBug06_PersistFailureNotEnqueued(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "f")
	if err := os.WriteFile(bad, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	q, err := worklease.New(worklease.Options{
		Clock: clock.NewFake(time.Now()), PersistPath: filepath.Join(bad, "q.json"),
	})
	if err != nil {
		t.Fatal(err)
	}
	defer q.Close()
	err = q.Enqueue(worklease.EnqueueSpec{ID: "6", Payload: []byte("p")})
	if err == nil {
		t.Fatal("expected persist error")
	}
	if len(q.List()) != 0 {
		t.Fatalf("job remained: %d", len(q.List()))
	}
}
