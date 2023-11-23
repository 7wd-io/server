package main

import (
	"7wd.io/adapter/analyst"
	"7wd.io/adapter/clock"
	"7wd.io/adapter/password"
	"7wd.io/adapter/playagain"
	"7wd.io/adapter/pusher"
	"7wd.io/adapter/repo"
	"7wd.io/adapter/token"
	"7wd.io/adapter/uuidf"
	"7wd.io/config"
	"7wd.io/domain"
	"7wd.io/infra/cent"
	"7wd.io/infra/pg"
	"7wd.io/infra/rds"
	"context"
	"fmt"
	swde "github.com/7wd-io/engine"
	"log/slog"
	"time"
)

func main() {
	ctx := context.Background()

	rdsc := rds.MustNew()
	pgc := pg.MustNew(ctx)
	cfgo := cent.New()
	psh := pusher.New(cfgo)

	roomRepo := repo.NewRoom(rdsc)
	gameRepo := repo.NewGame(pgc)
	userRepo := repo.NewUser(pgc)
	gameClockRepo := repo.NewGameClock(rdsc)
	sessionRepo := repo.NewSession(rdsc)
	playAgainStore := playagain.New(rdsc)

	anl := analyst.New(rdsc, pgc)

	accountSvc := domain.NewAccountService(
		userRepo,
		password.New(),
		clock.New(),
		token.New(config.C.Secret),
		uuidf.New(),
		sessionRepo,
		anl,
	)

	for {
		rooms, err := roomRepo.FindAll(ctx)

		if err != nil {
			slog.Error(err.Error())
			continue
		}

		for _, room := range rooms {
			// skip not started games
			if room.GameId == 0 {
				continue
			}

			gc, err := gameClockRepo.Find(ctx, room.GameId)

			if err != nil {
				slog.Error(fmt.Sprintf("game clock: %s (game id=%d)", err, room.GameId))
				continue
			}

			now := time.Now()

			timePassed := domain.TimeBank(now.Sub(gc.LastMoveAt))

			// skip live games
			if timePassed < gc.Values[gc.Turn] {
				continue
			}

			gc.Values[gc.Turn] = 0

			game, err := gameRepo.Find(ctx, domain.WithGameId(room.GameId))

			if err != nil {
				slog.Error(fmt.Sprintf("game clock: failed during gameRepo.Find: %w", err))
				continue
			}

			move := swde.NewMoveOver(swde.Nickname(gc.Turn), swde.Timeout)

			state, err := game.Move(gc.Turn, move)

			if err != nil {
				slog.Error(fmt.Sprintf("game clock: failed during moveOver: %w", err))
				continue
			}

			result := game.Over(state, now)

			if err := gameRepo.Update(ctx, game); err != nil {
				slog.Error(fmt.Sprintf("game clock: failed during update game: %w", err))
				continue
			}

			err = accountSvc.OnGameOver(ctx, domain.GameOverPayload{
				Game:    *game,
				Result:  result,
				Options: room.Options,
			})

			if err != nil {
				slog.Error(fmt.Sprintf("game clock: accountSvc.OnGameOver failed: %w", err))
			}

			go func() {
				err = psh.Publish(
					ctx,
					domain.ChGameUpdate(game.Id),
					domain.GameUpdatedPayload{
						Id:       game.Id,
						State:    state,
						Clock:    gc,
						LastMove: game.Log[len(game.Log)-1],
					},
				)

				if err != nil {
					slog.Error(err.Error())
				}
			}()

			if _, err = roomRepo.Delete(ctx, room.Id); err != nil {
				slog.Error(fmt.Sprintf("game clock: failed during delete room: %w", err))
				continue
			}

			go func() {
				err = psh.Publish(
					ctx,
					domain.ChRoomDelete,
					struct {
						Id domain.RoomId `json:"id"`
					}{
						Id: room.Id,
					},
				)

				if err != nil {
					slog.Error(err.Error())
				}
			}()

			if err = gameClockRepo.Delete(ctx, room.GameId); err != nil {
				slog.Error(fmt.Sprintf("game clock: failed during delete clock: %w", err))
				continue
			}

			if err := playAgainStore.Create(ctx, *game, room.Options); err != nil {
				slog.Error(fmt.Sprintf("game clock: cant create playAgain: %w", err))
				continue
			}
		}

		time.Sleep(time.Second)
	}
}
