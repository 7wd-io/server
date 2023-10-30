package domain

import (
	"context"
	"encoding/base64"
	"slices"
	"time"
)

const (
	RoomTtl = time.Minute * 60
)

type RoomId string

type Room struct {
	Id      RoomId      `json:"id"`
	Users   []UserId    `json:"users"`
	Options RoomOptions `json:"options"`
}

func (dst *Room) Join(u UserId) error {
	if slices.Contains(dst.Users, u) {
		return ErrAlreadyJoined
	}

	if len(dst.Users) >= dst.Options.Size {
		return ErrRoomIsFull
	}

	dst.Users = append(dst.Users, u)

	return nil
}

func (dst *Room) Leave(u UserId) {
	dst.Users = slices.DeleteFunc(dst.Users, func(id UserId) bool {
		return id == u
	})
}

func (dst *Room) Empty() bool {
	return len(dst.Users) == 0
}

type RoomOptions struct {
	// max 5
	Size int `json:"size"`
}

func NewRoomService(
	roomRepo RoomRepo,
	uuidf Uuidf,
) RoomService {
	return RoomService{
		roomRepo: roomRepo,
		uuidf:    uuidf,
	}
}

type RoomService struct {
	roomRepo RoomRepo
	uuidf    Uuidf
}

func (dst RoomService) List(ctx context.Context) ([]*Room, error) {
	return dst.roomRepo.Find()
}

func (dst RoomService) Create(ctx context.Context, pass Passport, o RoomOptions) (*Room, error) {
	room := &Room{
		Id:      RoomId(base64.RawURLEncoding.EncodeToString([]byte(dst.uuidf.Uuid().String()))),
		Users:   make([]UserId, 0, o.Size),
		Options: o,
	}

	if err := room.Join(pass.Id); err != nil {
		return nil, err
	}

	// @TODO проверить что у чела больше нет комнат
	// @TODO пуш обновы

	if err := dst.roomRepo.Save(ctx, room); err != nil {
		return nil, err
	}

	return room, nil
}

func (dst RoomService) Join(ctx context.Context, pass Passport, id RoomId) error {
	// @TODO lock
	room, err := dst.roomRepo.Find(ctx, id)

	if err != nil {
		return err
	}

	// @TODO пуш обновы

	if err := room.Join(pass.Id); err != nil {
		return err
	}

	return dst.roomRepo.Save(ctx, room)
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
