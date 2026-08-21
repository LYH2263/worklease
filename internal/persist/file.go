package persist

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"example.com/worklease/internal/job"
)

type Snapshot struct {
	Pending  []job.Job          `json:"pending"`
	Inflight map[string]job.Job `json:"inflight"`
}

type Store struct {
	mu   sync.Mutex
	path string
	f    *os.File
}

func New(path string) *Store { return &Store{path: path} }

func (s *Store) Save(snap Snapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	// Windows：重命名前必须先关掉已打开句柄
	if s.f != nil {
		_ = s.f.Close()
		s.f = nil
	}
	b, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, s.path); err != nil {
		return err
	}
	// 保存后重新打开，供 Close 释放；轮转测试依赖句柄
	f, err := os.OpenFile(s.path, os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	s.f = f
	return nil
}

func (s *Store) Load() (Snapshot, error) {
	b, err := os.ReadFile(s.path)
	if err != nil {
		return Snapshot{}, err
	}
	var snap Snapshot
	if err := json.Unmarshal(b, &snap); err != nil {
		return Snapshot{}, err
	}
	if snap.Inflight == nil {
		snap.Inflight = map[string]job.Job{}
	}
	return snap, nil
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.f = nil
	return nil
}

func (s *Store) Path() string { return s.path }
