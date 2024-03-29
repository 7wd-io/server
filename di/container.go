package di

import (
	"7wd.io/adapter/analyst"
	"7wd.io/adapter/bot"
	"7wd.io/adapter/clock"
	"7wd.io/adapter/dispatcher"
	"7wd.io/adapter/onliner"
	"7wd.io/adapter/password"
	"7wd.io/adapter/playagain"
	"7wd.io/adapter/repo"
	"7wd.io/adapter/token"
	"7wd.io/adapter/tx"
	"7wd.io/adapter/uuidf"
	"7wd.io/config"
	"7wd.io/infra/cent"
	"7wd.io/infra/pg"
	"7wd.io/infra/rds"
	"context"
	"github.com/centrifugal/gocent/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func MustNew() *C {
	pgc := pg.MustNew(context.Background())
	centfugo := cent.New()
	rdsc := rds.MustNew()

	return &C{
		Client: Client{
			Pg:   pgc,
			Rds:  rdsc,
			Cent: centfugo,
		},
		Repo: Repo{
			User:      repo.NewUser(pgc),
			Session:   repo.NewSession(rdsc),
			Room:      repo.NewRoom(rdsc),
			Game:      repo.NewGame(pgc),
			GameClock: repo.NewGameClock(rdsc),
		},

		Clock:        clock.New(),
		Tx:           tx.New(pgc),
		TokenFactory: token.New(config.C.Secret),
		UUIDFactory:  uuidf.New(),
		Pass:         password.New(),
		Onliner:      onliner.New(centfugo),
		Dispatcher:   dispatcher.New(),
		Bot:          bot.New(config.C.Bot.Endpoint),
		Analyst:      analyst.New(rdsc, pgc),
		PlayAgain:    playagain.New(rdsc),
	}
}

type C struct {
	Client Client
	Repo   Repo

	Clock        clock.Clock
	Tx           *tx.Tx
	TokenFactory token.F
	UUIDFactory  uuidf.F
	Pass         *password.Manager
	Onliner      *onliner.O
	Dispatcher   *dispatcher.D
	Bot          bot.B
	Analyst      analyst.A
	PlayAgain    playagain.PA
}

type Repo struct {
	User      repo.UserRepo
	Session   repo.SessionRepo
	Room      repo.RoomRepo
	Game      repo.GameRepo
	GameClock repo.GameClockRepo
}

type Client struct {
	Pg   *pgxpool.Pool
	Rds  *redis.Client
	Cent *gocent.Client
}
