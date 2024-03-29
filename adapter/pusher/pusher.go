package pusher

import (
	"7wd.io/domain"
	"context"
	"encoding/json"
	"github.com/centrifugal/gocent/v3"
)

func New(cent *gocent.Client) P {
	return P{
		cent: cent,
	}
}

type P struct {
	cent *gocent.Client
}

func (dst P) Publish(ctx context.Context, channel string, data interface{}) error {
	msg, _ := json.Marshal(data)

	_, err := dst.cent.Publish(ctx, channel, msg)

	return err
}

func (dst P) OnRoomCreated(ctx context.Context, payload interface{}) error {
	p, _ := payload.(domain.RoomCreatedPayload)

	return dst.Publish(ctx, domain.ChRoomCreate, p.Room)
}

func (dst P) OnRoomDeleted(ctx context.Context, payload interface{}) error {
	p, _ := payload.(domain.RoomDeletedPayload)

	return dst.Publish(ctx, domain.ChRoomDelete, struct {
		Id domain.RoomId `json:"id"`
	}{
		Id: p.Room.Id,
	})
}

func (dst P) OnRoomUpdated(ctx context.Context, payload interface{}) error {
	p, _ := payload.(domain.RoomUpdatedPayload)

	return dst.Publish(ctx, domain.ChRoomUpdate, p.Room)
}

func (dst P) OnGameUpdated(ctx context.Context, payload interface{}) error {
	p, _ := payload.(domain.GameUpdatedPayload)

	return dst.Publish(ctx, domain.ChGameUpdate(p.Id), payload)
}

func (dst P) OnPlayAgainUpdated(ctx context.Context, payload interface{}) error {
	p, _ := payload.(domain.PlayAgainUpdatedPayload)

	return dst.Publish(ctx, domain.ChPlayAgainUpdate(p.Game), payload)
}

func (dst P) OnPlayAgainApproved(ctx context.Context, payload interface{}) error {
	p, _ := payload.(domain.PlayAgainApprovedPayload)

	return dst.Publish(ctx, domain.ChPlayAgainApprove(p.Id), payload)
}
