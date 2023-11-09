package main

import (
	"7wd.io/config"
	"7wd.io/di"
	"7wd.io/domain"
	srv "7wd.io/http"
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
		c.Pusher,
		c.Dispatcher,
	)

	accountSvc := domain.NewAccountService(
		c.Repo.User,
		c.Pass,
		c.Clock,
		c.TokenFactory,
		c.UUIDFactory,
		c.Repo.Session,
	)

	roomSvc := domain.NewRoomService(
		c.Repo.Room,
		c.Repo.User,
		c.UUIDFactory,
		c.Dispatcher,
	)

	c.Dispatcher.
		On(
			domain.EventGameCreated,
			botSvc.OnGameCreated,
		).
		On(domain.EventGameUpdated).
		On(
			domain.EventGameOver,
			roomSvc.OnGameOver,
		).
		On(
			domain.EventAfterGameMove,
			botSvc.OnAfterGameMove,
		).
		On(
			domain.EventBotIsReadyToMove,
			gameSvc.OnEventBotIsReadyToMove,
		).
		On(domain.EventRoomCreated).
		On(domain.EventRoomUpdated).
		On(domain.EventRoomDeleted).
		On(domain.EventOnlineUpdated).
		On(domain.EventPlayAgainUpdated).
		On(domain.EventPlayAgainApproved)

	app.NewAccount(accountSvc).Bind(app)
	app.NewRoom(roomSvc).Bind(app)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", config.C.Port)))
}
