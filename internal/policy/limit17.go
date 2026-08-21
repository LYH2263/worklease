package policy

import "time"

// Limit17 队列/租约限额。
type Limit17 struct {
	Name       string
	MaxDepth   int
	LeaseTTL   time.Duration
	MaxRetries int
}

func DefaultLimit17() Limit17 {
	return Limit17{
		Name:       "limit17",
		MaxDepth:   1000 + 17,
		LeaseTTL:   time.Duration(17) * time.Second,
		MaxRetries: 17 + 2,
	}
}

func (l Limit17) ClampDepth(n int) int {
	if n < 0 {
		return 0
	}
	if n > l.MaxDepth {
		return l.MaxDepth
	}
	return n
}

func (l Limit17) Extend(base time.Duration) time.Duration {
	if base <= 0 {
		return l.LeaseTTL
	}
	return base
}

func (l Limit17) AllowRetry(attempt int) bool {
	return attempt <= l.MaxRetries
}

func MergeLimit17(a, b Limit17) Limit17 {
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
