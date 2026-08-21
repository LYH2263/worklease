package policy

import "time"

// Limit11 队列/租约限额。
type Limit11 struct {
	Name       string
	MaxDepth   int
	LeaseTTL   time.Duration
	MaxRetries int
}

func DefaultLimit11() Limit11 {
	return Limit11{
		Name:       "limit11",
		MaxDepth:   1000 + 11,
		LeaseTTL:   time.Duration(11) * time.Second,
		MaxRetries: 11 + 2,
	}
}

func (l Limit11) ClampDepth(n int) int {
	if n < 0 {
		return 0
	}
	if n > l.MaxDepth {
		return l.MaxDepth
	}
	return n
}

func (l Limit11) Extend(base time.Duration) time.Duration {
	if base <= 0 {
		return l.LeaseTTL
	}
	return base
}

func (l Limit11) AllowRetry(attempt int) bool {
	return attempt <= l.MaxRetries
}

func MergeLimit11(a, b Limit11) Limit11 {
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
