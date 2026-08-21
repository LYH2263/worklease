package policy

import "time"

// Limit12 队列/租约限额。
type Limit12 struct {
	Name       string
	MaxDepth   int
	LeaseTTL   time.Duration
	MaxRetries int
}

func DefaultLimit12() Limit12 {
	return Limit12{
		Name:       "limit12",
		MaxDepth:   1000 + 12,
		LeaseTTL:   time.Duration(12) * time.Second,
		MaxRetries: 12 + 2,
	}
}

func (l Limit12) ClampDepth(n int) int {
	if n < 0 {
		return 0
	}
	if n > l.MaxDepth {
		return l.MaxDepth
	}
	return n
}

func (l Limit12) Extend(base time.Duration) time.Duration {
	if base <= 0 {
		return l.LeaseTTL
	}
	return base
}

func (l Limit12) AllowRetry(attempt int) bool {
	return attempt <= l.MaxRetries
}

func MergeLimit12(a, b Limit12) Limit12 {
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
