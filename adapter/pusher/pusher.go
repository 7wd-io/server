package pusher

import (
	"7wd.io/domain"
	"context"
	"encoding/json"
	"fmt"
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

func (dst P) OnRoomCreated(ctx context.Context, payload interface{}) error {
	p, _ := payload.(domain.RoomCreatedPayload)

	dst.publish(ctx, "new_room", p.Room)

	return nil
}

func (dst P) OnRoomDeleted(ctx context.Context, payload interface{}) error {
	p, _ := payload.(domain.RoomDeletedPayload)

	dst.publish(ctx, "del_room", struct {
		Host domain.Nickname `json:"host"`
	}{
		Host: p.Room.Host,
	})

	return nil
}

func (dst P) OnRoomUpdated(ctx context.Context, payload interface{}) error {
	p, _ := payload.(domain.RoomUpdatedPayload)

	dst.publish(ctx, "upd_room", p.Room)

	return nil
}

func (dst P) OnGameCreated(ctx context.Context, payload interface{}) error {
	p, _ := payload.(domain.GameCreatedPayload)

	fmt.Println(p)

	return nil
}

func (dst P) OnGameUpdated(ctx context.Context, payload interface{}) error {
	p, _ := payload.(domain.GameUpdatedPayload)

	fmt.Println(p)

	return nil

	//dst.publish(
	//	fmt.Sprintf("upd_game_%d", p.Game.Id),
	//	struct {
	//		State    *game.State    `json:"state"`
	//		Clock    *domain.GameClock    `json:"clock"`
	//		LastMove domain.GameLogRecord `json:"lastMove"`
	//	}{
	//		State:    msg.State,
	//		Clock:    msg.Clock,
	//		LastMove: msg.LastMove,
	//	},
	//)
}

func (dst P) OnOnlineUpdated(ctx context.Context, payload interface{}) error {
	p, _ := payload.(domain.OnlineUpdatedPayload)

	fmt.Println(p)

	return nil
}

func (dst P) OnPlayAgainUpdated(ctx context.Context, payload interface{}) error {
	p, _ := payload.(domain.PlayAgainUpdatedPayload)

	return dst.publish(ctx, fmt.Sprintf("upd_play_again_%d", p.Game), payload)
}

func (dst P) OnPlayAgainApproved(ctx context.Context, payload interface{}) error {
	p, _ := payload.(domain.PlayAgainApprovedPayload)

	return dst.publish(ctx, fmt.Sprintf("play_again_approved_%d", p.Id), payload)
}

func (dst P) publish(ctx context.Context, channel string, data interface{}) error {
	msg, _ := json.Marshal(data)

	_, err := dst.cent.Publish(ctx, channel, msg)

	return err
}
