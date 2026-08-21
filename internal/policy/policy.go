package policy

import (
	"context"
	"time"
)

type Policy struct {
	LeaseTTL time.Duration
	MaxDepth int
}

func Default() Policy {
	return Policy{LeaseTTL: 2 * time.Second, MaxDepth: 1024}
}

func (p Policy) WaitClaim(ctx context.Context) error {
	_ = ctx
	return nil
}

func (p Policy) WaitHeartbeat(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
