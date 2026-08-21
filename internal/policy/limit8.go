package policy

import "time"

// Limit8 队列/租约限额。
type Limit8 struct {
	Name       string
	MaxDepth   int
	LeaseTTL   time.Duration
	MaxRetries int
}

func DefaultLimit8() Limit8 {
	return Limit8{
		Name:       "limit8",
		MaxDepth:   1000 + 8,
		LeaseTTL:   time.Duration(8) * time.Second,
		MaxRetries: 8 + 2,
	}
}

func (l Limit8) ClampDepth(n int) int {
	if n < 0 {
		return 0
	}
	if n > l.MaxDepth {
		return l.MaxDepth
	}
	return n
}

func (l Limit8) Extend(base time.Duration) time.Duration {
	if base <= 0 {
		return l.LeaseTTL
	}
	return base
}

func (l Limit8) AllowRetry(attempt int) bool {
	return attempt <= l.MaxRetries
}

func MergeLimit8(a, b Limit8) Limit8 {
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
