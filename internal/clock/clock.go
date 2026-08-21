package clock

import "time"

type Clock interface {
	Now() time.Time
}

type Real struct{}

func (Real) Now() time.Time { return time.Now() }

type Fake struct{ T time.Time }

func NewFake(t time.Time) *Fake { return &Fake{T: t} }

func (f *Fake) Now() time.Time { return f.T }

func (f *Fake) Advance(d time.Duration) {
	if f == nil {
		return
	}
	f.T = f.T.Add(d)
}
