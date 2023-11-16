package playagain

import (
	"7wd.io/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

func New(c *redis.Client) PA {
	return PA{
		rds: c,
	}
}

type PA struct {
	rds *redis.Client
}

func (dst PA) Create(ctx context.Context, game domain.Game, o domain.RoomOptions) error {
	return dst.setValue(
		ctx,
		dst.k(game.Id),
		domain.PlayAgainAgreement{
			Answers: map[domain.Nickname]*bool{
				game.HostNickname:  nil,
				game.GuestNickname: nil,
			},
			Options: o,
		},
		domain.PlayAgainWaiting+(time.Second*10),
	)
}

func (dst PA) Update(ctx context.Context, id domain.GameId, u domain.Nickname, value bool) (*domain.PlayAgainAgreement, error) {
	k := dst.k(id)
	found, err := dst.rds.Exists(ctx, k).Result()

	if err != nil {
		return nil, err
	}

	if found == 0 {
		return nil, errors.New("play again not available")
	}

	v := new(domain.PlayAgainAgreement)

	if err = dst.getValue(ctx, k, &v); err != nil {
		return nil, err
	}

	v.Answers[u] = &value

	if err = dst.setValue(ctx, k, v, redis.KeepTTL); err != nil {
		return nil, err
	}

	return v, nil
}

func (dst PA) k(id domain.GameId) string {
	return fmt.Sprintf("game:%d:playagain", id)
}

func (dst PA) setValue(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	v, err := json.Marshal(value)

	if err != nil {
		return nil
	}

	return dst.rds.Set(ctx, key, v, ttl).Err()
}

func (dst PA) getValue(ctx context.Context, key string, dest interface{}) error {
	v, err := dst.rds.Get(ctx, key).Bytes()

	if err != nil {
		return err
	}

	return json.Unmarshal(v, dest)
}
