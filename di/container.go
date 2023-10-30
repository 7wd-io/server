package di

import (
	"7wd.io/adapter/clock"
	"7wd.io/adapter/password"
	"7wd.io/adapter/repo"
	"7wd.io/adapter/token"
	"7wd.io/adapter/tx"
	"7wd.io/adapter/uuidf"
	"7wd.io/config"
	"7wd.io/infra/pg"
	"7wd.io/infra/rds"
	"context"
)

func MustNew() *C {
	pgc := pg.MustNew(context.Background())
	rdsc := rds.MustNew()

	return &C{
		Repo: Repo{
			User:    repo.NewUser(pgc),
			Session: repo.NewSession(rdsc),
			Room:    repo.NewRoom(rdsc),
		},

		Clock:        clock.New(),
		Tx:           tx.New(pgc),
		TokenFactory: token.New(config.C.Secret),
		UUIDFactory:  uuidf.New(),
		Pass:         password.New(),
	}
}

type C struct {
	Repo Repo

	Clock        clock.Clock
	Tx           *tx.Tx
	TokenFactory token.F
	UUIDFactory  uuidf.F
	Pass         *password.Manager
}

type Repo struct {
	User    repo.UserRepo
	Session repo.SessionRepo
	Room    repo.RoomRepo
}
