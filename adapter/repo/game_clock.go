package repo

import (
	"7wd.io/adapter/repo/internal/rds"
	"7wd.io/domain"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

func NewGameClock(c *redis.Client) GameClockRepo {
	return GameClockRepo{
		R: rds.R{Rds: c},
	}
}

type GameClockRepo struct {
	rds.R
}

func (dst GameClockRepo) Save(ctx context.Context, c *domain.GameClock) error {
	// game time ~ max 10-20 min
	return dst.Set(ctx, dst.k(c.Id), c, time.Hour*24)
}

func (dst GameClockRepo) Find(ctx context.Context, id domain.GameId) (*domain.GameClock, error) {
	c := new(domain.GameClock)

	if err := dst.Get(ctx, dst.k(id), c); err != nil {
		if err == redis.Nil {
			return nil, domain.ErrGameClockNotFound
		}

		return nil, err
	}

	return c, nil
}

func (dst GameClockRepo) Delete(ctx context.Context, id domain.GameId) error {
	return dst.Rds.Del(ctx, dst.k(id)).Err()
}

func (dst GameClockRepo) k(id domain.GameId) string {
	return fmt.Sprintf("clock:%d", id)
}
