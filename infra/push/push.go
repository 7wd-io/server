package push

import (
	"7wd.io/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/centrifugal/gocent/v3"
	"log/slog"
)

// @TODO mapping to own structs

func New() P {
	return P{}
}

type P struct {
	cent *gocent.Client
}

func (dst P) onRoomCreated(ctx context.Context, payload interface{}) {
	p, _ := payload.(domain.RoomCreatedPayload)

	dst.publish(ctx, "new_room", p.Room)
}

func (dst P) onRoomDeleted(ctx context.Context, payload interface{}) {
	p, _ := payload.(domain.RoomDeletedPayload)

	dst.publish(ctx, "del_room", struct {
		Host domain.Nickname `json:"host"`
	}{
		Host: p.Room.Host,
	})
}

func (dst P) onRoomUpdated(ctx context.Context, payload interface{}) {
	p, _ := payload.(domain.RoomUpdatedPayload)

	dst.publish(ctx, "upd_room", p.Room)
}

func (dst P) onGameCreated(ctx context.Context, payload interface{}) {
	p, _ := payload.(domain.GameCreatedPayload)

	fmt.Println(p)
}

func (dst P) onGameUpdated(ctx context.Context, payload interface{}) {
	p, _ := payload.(domain.GameUpdatedPayload)

	fmt.Println(p)

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

func (dst P) onOnlineUpdated(ctx context.Context, payload interface{}) {
	p, _ := payload.(domain.OnlineUpdatedPayload)

	fmt.Println(p)
}

func (dst P) onPlayAgainUpdated(ctx context.Context, payload interface{}) {
	p, _ := payload.(domain.PlayAgainUpdatedPayload)

	fmt.Println(p)
}

func (dst P) onPlayAgainApproved(ctx context.Context, payload interface{}) {
	p, _ := payload.(domain.PlayAgainApprovedPayload)

	fmt.Println(p)
}

func (dst P) publish(ctx context.Context, channel string, data interface{}) {
	msg, _ := json.Marshal(data)

	_, err := dst.cent.Publish(ctx, channel, msg)

	if err != nil {
		slog.Error("push.publish", slog.String("err", err.Error()))
	}
}
