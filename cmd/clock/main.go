package main

import (
	"7wd.io/adapter/analyst"
	"7wd.io/adapter/clock"
	"7wd.io/adapter/password"
	"7wd.io/adapter/playagain"
	"7wd.io/adapter/pusher"
	"7wd.io/adapter/repo"
	"7wd.io/adapter/token"
	txx "7wd.io/adapter/tx"
	"7wd.io/adapter/uuidf"
	"7wd.io/config"
	"7wd.io/domain"
	"7wd.io/infra/cent"
	"7wd.io/infra/pg"
	"7wd.io/infra/rds"
	"context"
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
	txer := txx.New(pgc)

	anl := analyst.New(rdsc, pgc)

	accountSvc := domain.NewAccountService(
		userRepo,
		password.New(),
		clock.New(),
		token.New(config.C.Secret),
		uuidf.New(),
		sessionRepo,
		anl,
		txer,
	)

	handleRoom := func(room *domain.Room) {
		// skip not started games
		if room.GameId == 0 {
			return
		}

		gc, err := gameClockRepo.Find(ctx, room.GameId)

		if err != nil {
			slog.Error(
				"game clock: gameClockRepo.Find fail",
				slog.Int("game id", int(room.GameId)),
				slog.String("err", err.Error()),
			)
			return
		}

		now := time.Now()

		timePassed := domain.TimeBank(now.Sub(gc.LastMoveAt))

		// skip live games
		if timePassed < gc.Values[gc.Turn] {
			return
		}

		gc.Values[gc.Turn] = 0

		tx, err := txer.Tx(ctx)

		if err != nil {
			slog.Error(
				"game clock: txer.Tx fail",
				slog.Int("game id", int(room.GameId)),
				slog.String("err", err.Error()),
			)

			return
		}

		defer func() {
			if errTx := tx.Rollback(ctx); errTx != nil {
				slog.Error(
					"game clock: tx.Rollback fail",
					slog.Int("game id", int(room.GameId)),
					slog.String("err", errTx.Error()),
				)
			}
		}()

		game, err := gameRepo.Find(
			ctx,
			domain.WithGameId(room.GameId),
			domain.WithGameTx(tx),
			domain.WithGameLock(),
		)

		if err != nil {
			slog.Error(
				"game clock: gameRepo.Find fail",
				slog.Int("game id", int(room.GameId)),
				slog.String("err", err.Error()),
			)

			return
		}

		move := swde.NewMoveOver(swde.Nickname(gc.Turn), swde.Timeout)

		state, err := game.Move(gc.Turn, move)

		if err != nil {
			slog.Error(
				"game clock: moveOver fail",
				slog.Int("game id", int(room.GameId)),
				slog.String("err", err.Error()),
			)

			return
		}

		result := game.Over(state, now)

		if err := gameRepo.Update(ctx, game, domain.WithGameTx(tx)); err != nil {
			slog.Error(
				"game clock: gameRepo.Update fail",
				slog.Int("game id", int(room.GameId)),
				slog.String("err", err.Error()),
			)

			return
		}

		err = accountSvc.OnGameOver(ctx, domain.GameOverPayload{
			Game:    *game,
			Result:  result,
			Options: room.Options,
		})

		if err != nil {
			slog.Error(
				"game clock: accountSvc.OnGameOver fail",
				slog.Int("game id", int(room.GameId)),
				slog.String("err", err.Error()),
			)
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
				slog.Error(
					"game clock: psh.Publish GameUpdated fail",
					slog.Int("game id", int(room.GameId)),
					slog.String("err", err.Error()),
				)
			}
		}()

		if _, err = roomRepo.Delete(ctx, room.Id); err != nil {
			slog.Error(
				"game clock: roomRepo.Delete fail",
				slog.Int("game id", int(room.GameId)),
				slog.String("err", err.Error()),
			)

			return
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
				slog.Error(
					"game clock: psh.Publish RoomDelete fail",
					slog.Int("game id", int(room.GameId)),
					slog.String("err", err.Error()),
				)
			}
		}()

		if err = gameClockRepo.Delete(ctx, room.GameId); err != nil {
			slog.Error(
				"game clock: gameClockRepo.Delete fail",
				slog.Int("game id", int(room.GameId)),
				slog.String("err", err.Error()),
			)

			return
		}

		if err = playAgainStore.Create(ctx, *game, room.Options); err != nil {
			slog.Error(
				"game clock: playAgainStore.Create fail",
				slog.Int("game id", int(room.GameId)),
				slog.String("err", err.Error()),
			)

			return
		}

		if err = tx.Commit(ctx); err != nil {
			slog.Error(
				"game clock: tx.Commit fail",
				slog.Int("game id", int(room.GameId)),
				slog.String("err", err.Error()),
			)

			return
		}
	}

	for {
		rooms, err := roomRepo.FindAll(ctx)

		if err != nil {
			slog.Error("game clock: roomRepo.FindAll fail", slog.String("err", err.Error()))
			continue
		}

		for _, room := range rooms {
			handleRoom(room)
		}

		time.Sleep(time.Second)
	}
}
