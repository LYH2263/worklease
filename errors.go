package worklease

import "errors"

var (
	ErrClosed    = errors.New("worklease: closed")
	ErrNotFound  = errors.New("worklease: job not found")
	ErrNotOwner  = errors.New("worklease: not lease owner")
	ErrEmpty     = errors.New("worklease: queue empty")
	ErrInvalid   = errors.New("worklease: invalid argument")
	ErrPersist   = errors.New("worklease: persist failed")
	ErrCanceled  = errors.New("worklease: canceled")
	ErrNoClock   = errors.New("worklease: no clock")
	ErrExpired   = errors.New("worklease: lease expired")
)
