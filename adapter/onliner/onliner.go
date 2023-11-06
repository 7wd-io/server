package onliner

import (
	"7wd.io/domain"
	"github.com/centrifugal/gocent/v3"
)

func New(client *gocent.Client) *O {
	return &O{client: client}
}

type O struct {
	client *gocent.Client
}

func (dst *O) Online() []domain.Nickname {
	return nil
}
