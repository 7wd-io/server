package domain

import (
	"context"
	"errors"
	swde "github.com/7wd-io/engine"
	"log/slog"
	"time"
)

const BotNickname Nickname = "bot"
const BotPlayAgainDelay = 3 * time.Second
const BotMoveDelay = 2 * time.Second

var BotNicknames = []Nickname{BotNickname, "b0t"}

func NewBotService(
	bot Bot,
	dispatcher Dispatcher,
) BotService {
	return BotService{
		bot:        bot,
		dispatcher: dispatcher,
	}
}

type BotService struct {
	bot        Bot
	dispatcher Dispatcher
}

func (dst BotService) tryMove(ctx context.Context, g *Game) {
	if !g.IsOver() && g.State().Me.Name == swde.Nickname(BotNickname) {
		go dst.move(ctx, g)
	}
}

func (dst BotService) move(ctx context.Context, g *Game) {
	time.Sleep(BotMoveDelay)

	move, err := dst.bot.GetMove(g)

	if err != nil {
		slog.Error(
			"dst.bot.GetMove",
			slog.String("err", err.Error()),
			slog.Int("game_id", int(g.Id)),
		)

		return
	}

	dst.dispatcher.Dispatch(
		ctx,
		EventBotIsReadyToMove,
		BotIsReadyToMovePayload{
			Game: g.Id,
			Move: move,
		},
	)
}

func (dst BotService) OnGameCreated(ctx context.Context, payload interface{}) error {
	p, ok := payload.(GameCreatedPayload)

	if !ok {
		return errors.New("BotService !ok := payload.(GameCreatedPayload)")
	}

	dst.tryMove(ctx, p.Game)

	return nil
}

func (dst BotService) OnAfterGameMove(ctx context.Context, payload interface{}) error {
	p, ok := payload.(AfterGameMovePayload)

	if !ok {
		return errors.New("BotService !ok := payload.(AfterGameMovePayload)")
	}

	dst.tryMove(ctx, p.Game)

	return nil
}
