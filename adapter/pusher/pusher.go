package pusher

import "github.com/centrifugal/gocent/v3"

func New(client *gocent.Client) *P {
	return &P{client: client}
}

type P struct {
	client *gocent.Client
}

func (dst *P) Push(msg interface{}) {
	// @TODO implement me
}
