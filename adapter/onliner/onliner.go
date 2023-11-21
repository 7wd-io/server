package onliner

import (
	"7wd.io/domain"
	"context"
	"github.com/centrifugal/gocent/v3"
)

func New(client *gocent.Client) *O {
	return &O{client: client}
}

type O struct {
	client *gocent.Client
}

func (dst *O) Online(ctx context.Context) ([]domain.Nickname, error) {
	result, err := dst.client.Presence(ctx, domain.ChOnline)

	if err != nil {
		return nil, err
	}

	// reserve 1 cap for bot
	online := make([]domain.Nickname, len(result.Presence)+1)

	i := 0
	for _, info := range result.Presence {
		online[i] = domain.Nickname(info.User)
		i++
	}

	// bot always online
	online[i] = domain.BotNickname

	return online, nil
}
