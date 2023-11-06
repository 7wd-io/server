package owatcher

import (
	"7wd.io/domain"
	"github.com/centrifugal/gocent/v3"
)

func New(client *gocent.Client) *W {
	return &W{client: client}
}

type W struct {
	client *gocent.Client
}

func (dst *W) Online() []domain.Nickname {
	return nil
}
