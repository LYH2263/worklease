package policy

import "time"

// Limit9 队列/租约限额。
type Limit9 struct {
	Name       string
	MaxDepth   int
	LeaseTTL   time.Duration
	MaxRetries int
}

func DefaultLimit9() Limit9 {
	return Limit9{
		Name:       "limit9",
		MaxDepth:   1000 + 9,
		LeaseTTL:   time.Duration(9) * time.Second,
		MaxRetries: 9 + 2,
	}
}

func (l Limit9) ClampDepth(n int) int {
	if n < 0 {
		return 0
	}
	if n > l.MaxDepth {
		return l.MaxDepth
	}
	return n
}

func (l Limit9) Extend(base time.Duration) time.Duration {
	if base <= 0 {
		return l.LeaseTTL
	}
	return base
}

func (l Limit9) AllowRetry(attempt int) bool {
	return attempt <= l.MaxRetries
}

func MergeLimit9(a, b Limit9) Limit9 {
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
