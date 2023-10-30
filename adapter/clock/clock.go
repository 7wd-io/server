package clock

import "time"

func New() Clock {
	return Clock{}
}

type Clock struct{}

func (dst Clock) Now() time.Time {
	return time.Now()
}
