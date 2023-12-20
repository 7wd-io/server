package http

import (
	"7wd.io/di"
	"7wd.io/domain"
	pgsuite "7wd.io/tt/suite/pg"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"path"
	"testing"
)

func Test_game(t *testing.T) {
	suite.Run(t, new(gameSuite))
}

type gameSuite struct {
	suite.Suite
	pgs *pgsuite.S
	c   *di.C
	srv *fiber.App
}

func (dst *gameSuite) SetupSuite() {
	dst.pgs.SetupSuite()

	c := di.MustNew()
	dst.c = c

	srv := New()

	gameSvc := domain.NewGameService(
		c.Clock,
		c.Repo.Room,
		c.Repo.Game,
		c.Repo.GameClock,
		c.Repo.User,
		c.Dispatcher,
	)

	NewGame(gameSvc).Bind(srv)

	dst.srv = srv

	// создать сервер
	// привязать роуты
	// вызвать fiber.test()

	// mute
}

func (dst *gameSuite) TearDownSuite() {
	dst.pgs.TearDownSuite()
}

func (dst *gameSuite) SetupTest() {
	dst.pgs.SetupTest(pgsuite.Options{
		Path: path.Join("http", "fixtures"),
	})
}

func (dst *gameSuite) TearDownTest() {
	dst.pgs.TearDownTest()
	dst.c.Client.Rds.FlushDB(context.Background())
}

func (dst *gameSuite) Test_Game1() {
	// Реквесты:
	// 	- создать игру
	//  - сделать ходы до конца

	dst.srv.Test()
}
