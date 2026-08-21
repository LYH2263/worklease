package policy

import "time"

// Limit13 队列/租约限额。
type Limit13 struct {
	Name       string
	MaxDepth   int
	LeaseTTL   time.Duration
	MaxRetries int
}

func DefaultLimit13() Limit13 {
	return Limit13{
		Name:       "limit13",
		MaxDepth:   1000 + 13,
		LeaseTTL:   time.Duration(13) * time.Second,
		MaxRetries: 13 + 2,
	}
}

func (l Limit13) ClampDepth(n int) int {
	if n < 0 {
		return 0
	}
	if n > l.MaxDepth {
		return l.MaxDepth
	}
	return n
}

func (l Limit13) Extend(base time.Duration) time.Duration {
	if base <= 0 {
		return l.LeaseTTL
	}
	return base
}

func (l Limit13) AllowRetry(attempt int) bool {
	return attempt <= l.MaxRetries
}

func MergeLimit13(a, b Limit13) Limit13 {
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
