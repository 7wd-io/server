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
	pusher Pusher,
) RoomService {
	return RoomService{
		roomRepo: roomRepo,
		userRepo: userRepo,
		uuidf:    uuidf,
		pusher:   pusher,
	}
}

type RoomService struct {
	roomRepo RoomRepo
	userRepo UserRepo
	uuidf    Uuidf
	pusher   Pusher
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

	go dst.pusher.Push(RoomCreated{Room: room})

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

	go dst.pusher.Push(RoomDeleted{Host: room.Host})

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

	go dst.pusher.Push(RoomUpdated{Room: room})

	return nil
}

func (dst RoomService) Leave(ctx context.Context, pass Passport, id RoomId) error {
	room, err := dst.roomRepo.Find(ctx, id)

	if err != nil {
		return err
	}

	if room.GameId != 0 {
		return ErrCantLeaveInProgressRoom
	}

	if room.Guest != pass.Nickname {
		return ErrRoomPlayerNotFound
	}

	room.Guest = ""
	room.GuestRating = 0

	if err := dst.roomRepo.Save(ctx, room); err != nil {
		return err
	}

	go dst.pusher.Push(RoomUpdated{Room: room})

	return nil
}

func (dst RoomService) Kick(ctx context.Context, pass Passport, id RoomId) error {
	room, err := dst.roomRepo.Find(ctx, id)

	if err != nil {
		return err
	}

	if room.GameId != 0 {
		return ErrCantLeaveInProgressRoom
	}

	if room.Guest != pass.Nickname {
		return ErrRoomPlayerNotFound
	}

	room.Guest = ""
	room.GuestRating = 0

	if err := dst.roomRepo.Save(ctx, room); err != nil {
		return err
	}

	go dst.pusher.Push(RoomUpdated{Room: room})

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
