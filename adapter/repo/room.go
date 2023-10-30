package repo

import (
	"7wd.io/adapter/repo/internal/rds"
	"7wd.io/domain"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func NewRoom(c *redis.Client) RoomRepo {
	return RoomRepo{
		R: rds.R{Client: c},
	}
}

type RoomRepo struct {
	rds.R
}

func (dst RoomRepo) Save(ctx context.Context, r *domain.Room) error {
	return dst.Set(ctx, dst.k(r.Id), r, domain.RoomTtl)
}

func (dst RoomRepo) Delete(ctx context.Context, id domain.RoomId) (*domain.Room, error) {
	s, err := dst.Find(ctx, id)

	if err != nil {
		return nil, err
	}

	if s == nil {
		return nil, domain.ErrRoomNotFound
	}

	return s, dst.Client.Del(ctx, dst.k(s.Id)).Err()
}

func (dst RoomRepo) Find(ctx context.Context, id domain.RoomId) (*domain.Room, error) {
	r := new(domain.Room)

	err := dst.Get(ctx, dst.k(id), r)

	if err != nil {
		return nil, err
	}

	return r, nil
}

func (dst RoomRepo) k(id domain.RoomId) string {
	return fmt.Sprintf("room:%s", id)
}
