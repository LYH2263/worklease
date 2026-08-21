package worklease

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"example.com/worklease/internal/clock"
)

// failingPersistPath 返回一个 PersistPath，其父目录是一个已存在的普通文件，
// 因此 Store.Save 中的 os.MkdirAll 必然失败。跨平台稳定，用于注入持久化失败。
func failingPersistPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	blocker := filepath.Join(dir, "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0o644); err != nil {
		t.Fatalf("create blocker: %v", err)
	}
	return filepath.Join(blocker, "queue.json")
}

// TestEnqueueRollsBackOnPersistFailure 复现「幽灵任务」：PersistPath 指错时
// 持久化失败，Enqueue 必须返回错误，并且不得把任务留在内存 pending 里。
// 否则 List 能看到、同进程 Claim 能拿到脏活、重启又消失，造成偶发双处理。
func TestEnqueueRollsBackOnPersistFailure(t *testing.T) {
	q, err := New(Options{
		Clock:       clock.NewFake(time.Now()),
		PersistPath: failingPersistPath(t),
		MaxDepth:    10,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	if err := q.Enqueue(EnqueueSpec{ID: "ghost", Payload: []byte("p")}); err == nil {
		_ = q.Close()
		t.Fatalf("Enqueue: expected persist failure, got nil")
	}
	_ = q.Close()

	// 核心断言：任务不得泄漏到 List。
	for _, v := range q.List() {
		if v.ID == "ghost" {
			t.Fatalf("ghost job leaked into List after persist failure: %+v", v)
		}
	}

	// 同进程 Claim 不得拿到这条脏活。
	if _, err := q.ClaimContext(context.Background(), "w1"); err == nil {
		t.Fatalf("Claim: expected no job, got a ghost")
	}

	// 失败不应计入 enqueued 指标。
	if got := q.Stats()["enqueued"]; got != 0 {
		t.Fatalf("enqueued metric = %d, want 0", got)
	}
}

// TestEnqueuePersistsOnSuccess 对照组：持久化正常时任务应入队且可被 Claim。
func TestEnqueuePersistsOnSuccess(t *testing.T) {
	dir := t.TempDir()
	q, err := New(Options{
		Clock:       clock.NewFake(time.Now()),
		PersistPath: filepath.Join(dir, "queue.json"),
		MaxDepth:    10,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	if err := q.Enqueue(EnqueueSpec{ID: "j1", Payload: []byte("p")}); err != nil {
		t.Fatalf("Enqueue: %v", err)
	}
	if got := q.Stats()["enqueued"]; got != 1 {
		t.Fatalf("enqueued metric = %d, want 1", got)
	}

	jv, err := q.ClaimContext(context.Background(), "w1")
	if err != nil {
		_ = q.Close()
		t.Fatalf("Claim: %v", err)
	}
	_ = q.Close()
	if jv.ID != "j1" {
		t.Fatalf("claimed %q, want j1", jv.ID)
	}
}
