package policy

import "time"

// Limit7 队列/租约限额。
type Limit7 struct {
	Name       string
	MaxDepth   int
	LeaseTTL   time.Duration
	MaxRetries int
}

func DefaultLimit7() Limit7 {
	return Limit7{
		Name:       "limit7",
		MaxDepth:   1000 + 7,
		LeaseTTL:   time.Duration(7) * time.Second,
		MaxRetries: 7 + 2,
	}
}

func (l Limit7) ClampDepth(n int) int {
	if n < 0 {
		return 0
	}
	if n > l.MaxDepth {
		return l.MaxDepth
	}
	return n
}

func (l Limit7) Extend(base time.Duration) time.Duration {
	if base <= 0 {
		return l.LeaseTTL
	}
	return base
}

func (l Limit7) AllowRetry(attempt int) bool {
	return attempt <= l.MaxRetries
}

func MergeLimit7(a, b Limit7) Limit7 {
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
