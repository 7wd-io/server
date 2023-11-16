package domain

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

const PlayAgainWaiting = time.Second * 60

type PlayAgainAgreement struct {
	Answers map[Nickname]*bool `json:"answers"`
	Options RoomOptions        `json:"options"`
}

func NewPlayAgainService(store PlayAgainStore) PlayAgainService {
	return PlayAgainService{
		store: store,
	}
}

type PlayAgainService struct {
	store PlayAgainStore
}

func (dst PlayAgainService) Update(ctx context.Context, game Game, u Nickname, value bool) error {
	return nil
}

func (dst PlayAgainService) OnGameOver(ctx context.Context, payload interface{}) error {
	p, ok := payload.(GameOverPayload)

	if !ok {
		return errors.New("func (dst PlayAgainService) OnGameOver !ok := payload.(GameOverPayload)")
	}

	var err error

	if err = dst.store.Create(ctx, p.Game, p.Options); err != nil {
		return err
	}

	if p.Game.GuestNickname == BotNickname {
		go func() {
			time.Sleep(BotPlayAgainDelay)

			if err = dst.Update(ctx, p.Game, BotNickname, true); err != nil {
				slog.Error(fmt.Sprintf("bot agree play again %w", err))
			}
		}()
	}

	return nil
}
