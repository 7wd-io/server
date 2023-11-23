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

func NewPlayAgainService(
	store PlayAgainStore,
	dispatcher Dispatcher,
	userRepo UserRepo,
	roomRepo RoomRepo,
	game GameCreator,
	gameRepo GameRepo,
) PlayAgainService {
	return PlayAgainService{
		store:      store,
		dispatcher: dispatcher,
		userRepo:   userRepo,
		roomRepo:   roomRepo,
		game:       game,
		gameRepo:   gameRepo,
	}
}

type PlayAgainService struct {
	store      PlayAgainStore
	dispatcher Dispatcher
	userRepo   UserRepo
	roomRepo   RoomRepo
	game       GameCreator
	gameRepo   GameRepo
}

func (dst PlayAgainService) UpdateById(ctx context.Context, id GameId, u Nickname, value bool) error {
	game, err := dst.gameRepo.Find(ctx, WithGameId(id))

	if err != nil {
		return err
	}

	return dst.Update(ctx, *game, u, value)
}

func (dst PlayAgainService) Update(ctx context.Context, game Game, u Nickname, value bool) error {
	pag, err := dst.store.Update(ctx, game.Id, u, value)

	if err != nil {
		return err
	}

	dst.dispatcher.Dispatch(
		ctx,
		EventPlayAgainUpdated,
		PlayAgainUpdatedPayload{
			Game:   game.Id,
			User:   u,
			Answer: value,
		},
	)

	agreement := true

	for _, answer := range pag.Answers {
		if answer == nil || *answer != true {
			agreement = false
			break
		}
	}

	if agreement {
		host, err := dst.userRepo.Find(ctx, WithUserNickname(game.HostNickname))

		if err != nil {
			return err
		}

		guest, err := dst.userRepo.Find(ctx, WithUserNickname(game.GuestNickname))

		if err != nil {
			return err
		}

		// reset options that don't make sense
		if pag.Options.MinRating != 0 {
			pag.Options.MinRating = 0
		}

		nextGame, err := dst.game.Create(ctx, *host, *guest, pag.Options)

		if err != nil {
			return err
		}

		room := &Room{
			GameId:      nextGame.Id,
			Host:        host.Nickname,
			HostRating:  host.Rating,
			Guest:       guest.Nickname,
			GuestRating: guest.Rating,
			Options:     pag.Options,
		}

		if err = dst.roomRepo.Save(ctx, room); err != nil {
			return err
		}

		dst.dispatcher.Dispatch(
			ctx,
			EventRoomCreated,
			RoomCreatedPayload{Room: *room},
		)

		dst.dispatcher.Dispatch(
			ctx,
			EventPlayAgainApproved,
			PlayAgainApprovedPayload{
				Id:   game.Id,
				Next: nextGame.Id,
			},
		)
	}

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
