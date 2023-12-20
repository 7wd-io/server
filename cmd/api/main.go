package main

import (
	"7wd.io/adapter/pusher"
	"7wd.io/config"
	"7wd.io/di"
	"7wd.io/domain"
	srv "7wd.io/http"
	"7wd.io/infra/cent"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"log/slog"
	"os"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	c := di.MustNew()

	app := srv.New()

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong 👋")
	})

	app.Get("/secret", func(c *fiber.Ctx) error {
		return c.SendString("love")
	})

	botSvc := domain.NewBotService(c.Bot, c.Dispatcher)

	gameSvc := domain.NewGameService(
		c.Clock,
		c.Repo.Room,
		c.Repo.Game,
		c.Repo.GameClock,
		c.Repo.User,
		c.Dispatcher,
	)

	push := pusher.New(cent.New())

	accountSvc := domain.NewAccountService(
		c.Repo.User,
		c.Pass,
		c.Clock,
		c.TokenFactory,
		c.UUIDFactory,
		c.Repo.Session,
		c.Analyst,
	)

	roomSvc := domain.NewRoomService(
		c.Repo.Room,
		c.Repo.User,
		c.UUIDFactory,
		c.Dispatcher,
	)

	onlineSvc := domain.NewOnlineService(c.Onliner, c.Analyst)

	playAgainSvc := domain.NewPlayAgainService(
		c.PlayAgain,
		c.Dispatcher,
		c.Repo.User,
		c.Repo.Room,
		gameSvc,
		c.Repo.Game,
	)

	c.Dispatcher.
		On(
			domain.EventGameCreated,
			botSvc.OnGameCreated,
		).
		On(
			domain.EventGameUpdated,
			push.OnGameUpdated,
		).
		On(
			domain.EventGameOver,
			roomSvc.OnGameOver,
			accountSvc.OnGameOver,
			playAgainSvc.OnGameOver,
		).
		On(
			domain.EventAfterGameMove,
			botSvc.OnAfterGameMove,
		).
		On(
			domain.EventBotIsReadyToMove,
			gameSvc.OnEventBotIsReadyToMove,
		).
		On(
			domain.EventRoomCreated,
			push.OnRoomCreated,
		).
		On(
			domain.EventRoomUpdated,
			push.OnRoomUpdated,
		).
		On(
			domain.EventRoomDeleted,
			push.OnRoomDeleted,
		).
		On(
			domain.EventRoomStarted,
			gameSvc.OnRoomStarted,
		).
		On(
			domain.EventPlayAgainUpdated,
			push.OnPlayAgainUpdated,
		).
		On(
			domain.EventPlayAgainApproved,
			push.OnPlayAgainApproved,
		)

	srv.NewAccount(accountSvc).Bind(app)
	srv.NewRoom(roomSvc).Bind(app)
	srv.NewOnline(onlineSvc).Bind(app)
	srv.NewGame(gameSvc).Bind(app)
	srv.NewPlayAgain(playAgainSvc).Bind(app)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", config.C.Port)))
}
