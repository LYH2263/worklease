package policy

import "time"

// Limit10 队列/租约限额。
type Limit10 struct {
	Name       string
	MaxDepth   int
	LeaseTTL   time.Duration
	MaxRetries int
}

func DefaultLimit10() Limit10 {
	return Limit10{
		Name:       "limit10",
		MaxDepth:   1000 + 10,
		LeaseTTL:   time.Duration(10) * time.Second,
		MaxRetries: 10 + 2,
	}
}

func (l Limit10) ClampDepth(n int) int {
	if n < 0 {
		return 0
	}
	if n > l.MaxDepth {
		return l.MaxDepth
	}
	return n
}

func (l Limit10) Extend(base time.Duration) time.Duration {
	if base <= 0 {
		return l.LeaseTTL
	}
	return base
}

func (l Limit10) AllowRetry(attempt int) bool {
	return attempt <= l.MaxRetries
}

func MergeLimit10(a, b Limit10) Limit10 {
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
