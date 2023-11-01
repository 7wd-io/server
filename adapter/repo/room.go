package repo

import (
	"7wd.io/adapter/repo/internal/rds"
	"7wd.io/domain"
	"context"
	"fmt"
	"github.com/google/uuid"
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
	// @TODO transaction
	if err := dst.Client.SAdd(ctx, dst.keyList(), r.Id).Err(); err != nil {
		return err
	}

	return dst.Set(ctx, dst.keyItem(r.Id), r, domain.RoomTtl)
}

func (dst RoomRepo) Delete(ctx context.Context, id domain.RoomId) (*domain.Room, error) {
	if err := dst.Client.SRem(ctx, dst.keyList(), uuid.UUID(id).String()).Err(); err != nil {
		return nil, err
	}

	s, err := dst.Find(ctx, id)

	if err != nil {
		return nil, err
	}

	if s == nil {
		return nil, domain.ErrRoomNotFound
	}

	return s, dst.Client.Del(ctx, dst.keyItem(s.Id)).Err()
}

func (dst RoomRepo) Find(ctx context.Context, id domain.RoomId) (*domain.Room, error) {
	r := new(domain.Room)

	err := dst.Get(ctx, dst.keyItem(id), r)

	if err != nil {
		return nil, err
	}

	return r, nil
}

func (dst RoomRepo) FindAll(ctx context.Context) ([]*domain.Room, error) {
	members, err := dst.Client.SMembers(ctx, dst.keyList()).Result()

	if err != nil {
		return nil, err
	}

	rooms := make([]*domain.Room, len(members))

	for k, v := range members {
		room := new(domain.Room)

		if err = dst.Get(ctx, dst.keyItem(domain.RoomId(uuid.MustParse(v))), room); err != nil {
			return nil, err
		}

		rooms[k] = room
	}

	return rooms, nil
}

func (dst RoomRepo) keyItem(id domain.RoomId) string {
	return fmt.Sprintf("room:%s", id)
}

func (dst RoomRepo) keyList() string {
	return "rooms"
}
