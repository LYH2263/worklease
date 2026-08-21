package policy

import "time"

// Limit15 队列/租约限额。
type Limit15 struct {
	Name       string
	MaxDepth   int
	LeaseTTL   time.Duration
	MaxRetries int
}

func DefaultLimit15() Limit15 {
	return Limit15{
		Name:       "limit15",
		MaxDepth:   1000 + 15,
		LeaseTTL:   time.Duration(15) * time.Second,
		MaxRetries: 15 + 2,
	}
}

func (l Limit15) ClampDepth(n int) int {
	if n < 0 {
		return 0
	}
	if n > l.MaxDepth {
		return l.MaxDepth
	}
	return n
}

func (l Limit15) Extend(base time.Duration) time.Duration {
	if base <= 0 {
		return l.LeaseTTL
	}
	return base
}

func (l Limit15) AllowRetry(attempt int) bool {
	return attempt <= l.MaxRetries
}

func MergeLimit15(a, b Limit15) Limit15 {
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
