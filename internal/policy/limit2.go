package policy

import "time"

// Limit2 队列/租约限额。
type Limit2 struct {
	Name       string
	MaxDepth   int
	LeaseTTL   time.Duration
	MaxRetries int
}

func DefaultLimit2() Limit2 {
	return Limit2{
		Name:       "limit2",
		MaxDepth:   1000 + 2,
		LeaseTTL:   time.Duration(2) * time.Second,
		MaxRetries: 2 + 2,
	}
}

func (l Limit2) ClampDepth(n int) int {
	if n < 0 {
		return 0
	}
	if n > l.MaxDepth {
		return l.MaxDepth
	}
	return n
}

func (l Limit2) Extend(base time.Duration) time.Duration {
	if base <= 0 {
		return l.LeaseTTL
	}
	return base
}

func (l Limit2) AllowRetry(attempt int) bool {
	return attempt <= l.MaxRetries
}

func MergeLimit2(a, b Limit2) Limit2 {
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
