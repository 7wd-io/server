package domain

import (
	"context"
	"errors"
	"fmt"
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
	HostRating  Rating      `json:"hostRating"`
	Guest       Nickname    `json:"guest,omitempty"`
	GuestRating Rating      `json:"guestRating,omitempty"`
	Options     RoomOptions `json:"options"`
	// set if room converted to game
	GameId GameId `json:"gameId,omitempty"`
}

type RoomOptions struct {
	Fast         bool     `json:"fast,omitempty"`
	MinRating    Rating   `json:"minRating,omitempty" validate:"omitempty,max=2000"`
	Enemy        Nickname `json:"enemy,omitempty"`
	PromoWonders bool     `json:"promoWonders"`
	TimeBank     TimeBank `json:"timeBank,omitempty"`
}

func (dst RoomOptions) Clock() TimeBank {
	// replace Fast to exactly time bank in future
	if dst.TimeBank != 0 {
		return dst.TimeBank
	}

	if dst.Fast {
		return timeBankFast
	}

	return timeBankDefault
}

func NewRoomService(
	roomRepo RoomRepo,
	userRepo UserRepo,
	uuidf Uuidf,
	dispatcher Dispatcher,
) RoomService {
	return RoomService{
		roomRepo:   roomRepo,
		userRepo:   userRepo,
		uuidf:      uuidf,
		dispatcher: dispatcher,
	}
}

type RoomService struct {
	roomRepo   RoomRepo
	userRepo   UserRepo
	uuidf      Uuidf
	dispatcher Dispatcher
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

	dst.dispatcher.Dispatch(
		ctx,
		EventRoomCreated,
		RoomCreatedPayload{Room: room},
	)

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

	dst.dispatcher.Dispatch(
		ctx,
		EventRoomDeleted,
		RoomDeletedPayload{Room: room},
	)

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

	dst.dispatcher.Dispatch(
		ctx,
		EventRoomUpdated,
		RoomUpdatedPayload{Room: *room},
	)

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

	dst.dispatcher.Dispatch(
		ctx,
		EventRoomUpdated,
		RoomUpdatedPayload{Room: *room},
	)

	return nil
}

func (dst RoomService) Start(ctx context.Context, pass Passport, id RoomId) error {
	room, err := dst.roomRepo.Find(ctx, id)

	if err != nil {
		return err
	}

	if room.Host != pass.Nickname {
		return ErrActionNotAllowed
	}

	if room.Guest == "" {
		return ErrActionNotAllowed
	}

	host, err := dst.userRepo.Find(ctx, WithUserNickname(room.Host))

	if err != nil {
		return err
	}

	guest, err := dst.userRepo.Find(ctx, WithUserNickname(room.Guest))

	if err != nil {
		return err
	}

	dst.dispatcher.Dispatch(
		ctx,
		EventRoomStarted,
		RoomStartedPayload{
			Host:  *host,
			Guest: *guest,
			Room:  *room,
		},
	)

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

	dst.dispatcher.Dispatch(
		ctx,
		EventRoomUpdated,
		RoomUpdatedPayload{Room: *room},
	)

	return nil
}

func (dst RoomService) OnGameOver(ctx context.Context, payload interface{}) error {
	p, ok := payload.(GameOverPayload)

	if !ok {
		return errors.New("func (dst RoomService) OnGameOver !ok := payload.(GameOverPayload)")
	}

	room, err := dst.roomRepo.FindByGame(ctx, p.Game.Id)

	if err != nil {
		return fmt.Errorf("func (dst RoomService) OnGameOver: %w", err)
	}

	if room == nil {
		return nil
	}

	if _, err = dst.roomRepo.Delete(ctx, room.Id); err != nil {
		return fmt.Errorf("dst.roomRepo.Delete(ctx, room.Id): %w", err)
	}

	dst.dispatcher.Dispatch(
		ctx,
		EventRoomDeleted,
		RoomDeletedPayload{Room: room},
	)

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
