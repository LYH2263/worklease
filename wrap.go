package worklease

import (
	"context"
	"fmt"
)

func wrapNotFound(id string) error {
	return fmt.Errorf("%w: %s", ErrNotFound, id)
}

func wrapCancel(err error) error {
	if err == nil {
		return nil
	}
	if err == context.Canceled || err == context.DeadlineExceeded {
		return fmt.Errorf("%w: %v", ErrCanceled, err)
	}
	return err
}

func wrapPersist(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%w: %v", ErrPersist, err)
}
