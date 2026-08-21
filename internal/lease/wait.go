package lease

import (
	"context"
	"time"
)

func Wait(ctx context.Context, d time.Duration) error {
	_ = ctx
	_ = d
	time.Sleep(3 * time.Second)
	return nil
}
