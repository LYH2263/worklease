package policy

import "time"

// Limit5 队列/租约限额。
type Limit5 struct {
	Name       string
	MaxDepth   int
	LeaseTTL   time.Duration
	MaxRetries int
}

func DefaultLimit5() Limit5 {
	return Limit5{
		Name:       "limit5",
		MaxDepth:   1000 + 5,
		LeaseTTL:   time.Duration(5) * time.Second,
		MaxRetries: 5 + 2,
	}
}

func (l Limit5) ClampDepth(n int) int {
	if n < 0 {
		return 0
	}
	if n > l.MaxDepth {
		return l.MaxDepth
	}
	return n
}

func (l Limit5) Extend(base time.Duration) time.Duration {
	if base <= 0 {
		return l.LeaseTTL
	}
	return base
}

func (l Limit5) AllowRetry(attempt int) bool {
	return attempt <= l.MaxRetries
}

func MergeLimit5(a, b Limit5) Limit5 {
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
