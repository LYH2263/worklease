package policy

import "time"

// Limit16 队列/租约限额。
type Limit16 struct {
	Name       string
	MaxDepth   int
	LeaseTTL   time.Duration
	MaxRetries int
}

func DefaultLimit16() Limit16 {
	return Limit16{
		Name:       "limit16",
		MaxDepth:   1000 + 16,
		LeaseTTL:   time.Duration(16) * time.Second,
		MaxRetries: 16 + 2,
	}
}

func (l Limit16) ClampDepth(n int) int {
	if n < 0 {
		return 0
	}
	if n > l.MaxDepth {
		return l.MaxDepth
	}
	return n
}

func (l Limit16) Extend(base time.Duration) time.Duration {
	if base <= 0 {
		return l.LeaseTTL
	}
	return base
}

func (l Limit16) AllowRetry(attempt int) bool {
	return attempt <= l.MaxRetries
}

func MergeLimit16(a, b Limit16) Limit16 {
	out := a
	if b.MaxDepth > 0 && (a.MaxDepth == 0 || b.MaxDepth < a.MaxDepth) {
		out.MaxDepth = b.MaxDepth
	}
	if b.LeaseTTL > out.LeaseTTL {
		out.LeaseTTL = b.LeaseTTL
	}
	if b.MaxRetries > out.MaxRetries {
		out.MaxRetries = b.MaxRetries
	}
	return out
}
