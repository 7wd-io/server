package main

import (
	"7wd.io/adapter/analyst"
	"7wd.io/adapter/onliner"
	"7wd.io/adapter/pusher"
	"7wd.io/adapter/repo"
	"7wd.io/domain"
	"7wd.io/infra/cent"
	"7wd.io/infra/pg"
	"7wd.io/infra/rds"
	"context"
	"log/slog"
	"os"
	"time"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
	ctx := context.Background()

	cfgo := cent.New()
	rdsc := rds.MustNew()
	pgc := pg.MustNew(ctx)
	psh := pusher.New(cfgo)

	roomRepo := repo.NewRoom(rdsc)
	anl := analyst.New(rdsc, pgc)
	watcher := onliner.New(cfgo)

	var err error
	var players []domain.Nickname
	var rooms []*domain.Room

	for {
		players, err = watcher.Online(ctx)

		if err != nil {
			slog.Error(err.Error())
			continue
		}

		rooms, err = roomRepo.FindAll(ctx)

		if err != nil {
			slog.Error(err.Error())
			continue
		}

		playersSearch := make(map[domain.Nickname]struct{}, len(players))

		for _, v := range players {
			playersSearch[v] = struct{}{}
		}

		for _, room := range rooms {
			_, isHostOnline := playersSearch[room.Host]

			// not empty gameId means game started and room must be show to support observe/play opportunity
			if !isHostOnline && room.GameId == 0 {
				_, err = roomRepo.Delete(ctx, room.Id)

				if err != nil {
					slog.Error(err.Error())
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

			}

			// kick if guest not online
			if room.GameId == 0 && room.Guest != "" {
				_, isGuestOnline := playersSearch[room.Guest]

				if !isGuestOnline {
					room.Guest = ""
					room.GuestRating = 0

					if err = roomRepo.Save(ctx, room); err != nil {
						slog.Error(err.Error())
						continue
					}

					go func() {
						err = psh.Publish(ctx, domain.ChRoomUpdate, room)

						if err != nil {
							slog.Error(err.Error())
						}
					}()
				}
			}
		}

		var online domain.UsersPreview

		if len(players) > 0 {
			online, err = anl.Ratings(ctx, players...)

			if err != nil {
				slog.Error(err.Error())
				continue
			}
		}

		go func() {
			err = psh.Publish(ctx, domain.ChOnline, online)

			if err != nil {
				slog.Error(err.Error())
			}
		}()

		time.Sleep(time.Second * 3)
	}
}
