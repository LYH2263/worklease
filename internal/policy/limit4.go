package policy

import "time"

// Limit4 队列/租约限额。
type Limit4 struct {
	Name       string
	MaxDepth   int
	LeaseTTL   time.Duration
	MaxRetries int
}

func DefaultLimit4() Limit4 {
	return Limit4{
		Name:       "limit4",
		MaxDepth:   1000 + 4,
		LeaseTTL:   time.Duration(4) * time.Second,
		MaxRetries: 4 + 2,
	}
}

func (l Limit4) ClampDepth(n int) int {
	if n < 0 {
		return 0
	}
	if n > l.MaxDepth {
		return l.MaxDepth
	}
	return n
}

func (l Limit4) Extend(base time.Duration) time.Duration {
	if base <= 0 {
		return l.LeaseTTL
	}
	return base
}

func (l Limit4) AllowRetry(attempt int) bool {
	return attempt <= l.MaxRetries
}

func MergeLimit4(a, b Limit4) Limit4 {
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
