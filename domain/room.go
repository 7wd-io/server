package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

const (
	RoomTtl = time.Minute * 60
)

type RoomId uuid.UUID

type Room struct {
	Id          RoomId      `json:"id"`
	Host        Nickname    `json:"host"`
	HostRating  int         `json:"hostRating"`
	Guest       Nickname    `json:"guest,omitempty"`
	GuestRating int         `json:"guestRating,omitempty"`
	Options     RoomOptions `json:"options"`
	// set if room converted to game
	GameId GameId `json:"gameId,omitempty"`
}

//func (dst *Room) Join(u UserId) error {
//	if slices.Contains(dst.Users, u) {
//		return ErrAlreadyJoined
//	}
//
//	if len(dst.Users) >= dst.Options.Size {
//		return ErrRoomIsFull
//	}
//
//	dst.Users = append(dst.Users, u)
//
//	return nil
//}
//
//func (dst *Room) Leave(u UserId) {
//	dst.Users = slices.DeleteFunc(dst.Users, func(id UserId) bool {
//		return id == u
//	})
//}

type RoomOptions struct {
	Fast         bool          `json:"fast,omitempty"`
	MinRating    int           `json:"minRating,omitempty" validate:"omitempty,max=2000"`
	Enemy        Nickname      `json:"enemy,omitempty"`
	PromoWonders bool          `json:"promoWonders"`
	TimeBank     time.Duration `json:"timeBank,omitempty"`
}

func NewRoomService(
	roomRepo RoomRepo,
	userRepo UserRepo,
	uuidf Uuidf,
) RoomService {
	return RoomService{
		roomRepo: roomRepo,
		userRepo: userRepo,
		uuidf:    uuidf,
	}
}

type RoomService struct {
	roomRepo RoomRepo
	userRepo UserRepo
	uuidf    Uuidf
}

func (dst RoomService) List(ctx context.Context) ([]*Room, error) {
	return dst.roomRepo.FindAll(ctx)
}

func (dst RoomService) Create(ctx context.Context, pass Passport, o RoomOptions) (*Room, error) {
	if o.Enemy != "" {
		enemy, err := dst.userRepo.Find(ctx, WithUserNickname(o.Enemy))

		if err != nil {
			return nil, err
		}

		if enemy.Nickname == pass.Nickname {
			return nil, ErrInvalidRoomOptions
		}
	}

	if err := dst.alreadyJoin(ctx, pass); err != nil {
		return nil, err
	}

	room := &Room{
		Id:         RoomId(dst.uuidf.Uuid()),
		Host:       pass.Nickname,
		HostRating: pass.Rating,
		Options:    o,
	}

	if err := dst.roomRepo.Save(ctx, room); err != nil {
		return nil, err
	}

	// @TODO cent

	return room, nil
}

func (dst RoomService) Delete(ctx context.Context, pass Passport, id RoomId) error {
	room, err := dst.roomRepo.Find(ctx, id)

	if err != nil {
		return err
	}

	if room.Host != pass.Nickname {
		return ErrOnlyHostCanRemoveRoom
	}

	if room.GameId != 0 {
		return ErrCantRemoveInProgressRoom
	}

	if _, err := dst.roomRepo.Delete(ctx, id); err != nil {
		return err
	}

	// @TODO cent

	return nil
}

func (dst RoomService) Join(ctx context.Context, pass Passport, id RoomId) error {
	if err := dst.alreadyJoin(ctx, pass); err != nil {
		return err
	}

	room, err := dst.roomRepo.Find(ctx, id)

	if err != nil {
		return err
	}

	if room.Guest != "" {
		return ErrRoomIsFull
	}

	if room.Options.Enemy != "" && pass.Nickname != room.Options.Enemy {
		return ErrJoinToTheRoomRestricted
	}

	if room.Options.MinRating != 0 && pass.Rating < room.Options.MinRating {
		return ErrJoinToTheRoomRestricted
	}

	room.Guest = pass.Nickname
	room.GuestRating = pass.Rating

	if err := dst.roomRepo.Save(ctx, room); err != nil {
		return err
	}

	// @TODO cent

	return nil
}

func (dst RoomService) Leave(ctx context.Context, pass Passport, id RoomId) error {
	room, err := dst.roomRepo.Find(ctx, id)

	if err != nil {
		return err
	}

	room.Leave(pass.Id)

	if room.Empty() {
		_, err = dst.roomRepo.Delete(ctx, id)

		if err != nil {
			return err
		}

		// @TODO пуш обновы
	} else {
		if err = dst.roomRepo.Save(ctx, room); err != nil {
			return err
		}

		// @TODO пуш обновы
	}

	return nil
}

func (dst RoomService) alreadyJoin(ctx context.Context, pass Passport) error {
	rooms, err := dst.roomRepo.FindAll(ctx)

	if err != nil {
		return err
	}

	for _, v := range rooms {
		if pass.Nickname == v.Host || pass.Nickname == v.Guest {
			return ErrOneRoomPerPlayer
		}
	}

	return nil
}
