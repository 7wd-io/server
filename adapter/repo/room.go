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
		R: rds.R{Rds: c},
	}
}

type RoomRepo struct {
	rds.R
}

func (dst RoomRepo) Save(ctx context.Context, r *domain.Room) error {
	// @TODO transaction
	if err := dst.Rds.SAdd(ctx, dst.keyList(), r.Id.String()).Err(); err != nil {
		return err
	}

	if r.GameId != 0 {
		if err := dst.Set(ctx, dst.keyGameRelation(r.GameId), r.Id, domain.RoomTtl); err != nil {
			return err
		}
	}

	return dst.Set(ctx, dst.keyItem(r.Id), r, domain.RoomTtl)
}

func (dst RoomRepo) Delete(ctx context.Context, id domain.RoomId) (*domain.Room, error) {
	if err := dst.Rds.SRem(ctx, dst.keyList(), id.String()).Err(); err != nil {
		return nil, err
	}

	s, err := dst.Find(ctx, id)

	if err != nil {
		return nil, err
	}

	if s == nil {
		return nil, domain.ErrRoomNotFound
	}

	// @TODO remove relation

	return s, dst.Rds.Del(ctx, dst.keyItem(s.Id)).Err()
}

func (dst RoomRepo) Find(ctx context.Context, id domain.RoomId) (*domain.Room, error) {
	r := new(domain.Room)

	err := dst.Get(ctx, dst.keyItem(id), r)

	if err != nil {
		return nil, err
	}

	return r, nil
}

func (dst RoomRepo) FindByGame(ctx context.Context, id domain.GameId) (*domain.Room, error) {
	var roomId domain.RoomId

	if err := dst.Get(ctx, dst.keyGameRelation(id), &roomId); err != nil {
		return nil, err
	}

	return dst.Find(ctx, roomId)
}

func (dst RoomRepo) FindAll(ctx context.Context) ([]*domain.Room, error) {
	members, err := dst.Rds.SMembers(ctx, dst.keyList()).Result()

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

func (dst RoomRepo) keyGameRelation(id domain.GameId) string {
	return fmt.Sprintf("room:game:%d", id)
}

func (dst RoomRepo) keyList() string {
	return "rooms"
}
