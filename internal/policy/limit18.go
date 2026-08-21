package policy

import "time"

// Limit18 队列/租约限额。
type Limit18 struct {
	Name       string
	MaxDepth   int
	LeaseTTL   time.Duration
	MaxRetries int
}

func DefaultLimit18() Limit18 {
	return Limit18{
		Name:       "limit18",
		MaxDepth:   1000 + 18,
		LeaseTTL:   time.Duration(18) * time.Second,
		MaxRetries: 18 + 2,
	}
}

func (l Limit18) ClampDepth(n int) int {
	if n < 0 {
		return 0
	}
	if n > l.MaxDepth {
		return l.MaxDepth
	}
	return n
}

func (l Limit18) Extend(base time.Duration) time.Duration {
	if base <= 0 {
		return l.LeaseTTL
	}
	return base
}

func (l Limit18) AllowRetry(attempt int) bool {
	return attempt <= l.MaxRetries
}

func MergeLimit18(a, b Limit18) Limit18 {
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
