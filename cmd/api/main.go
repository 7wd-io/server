package main

import (
	"7wd.io/config"
	"7wd.io/di"
	"7wd.io/domain"
	"7wd.io/infra/http"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"log/slog"
	"os"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	c := di.MustNew()

	app := http.NewApp()

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong 👋")
	})

	app.Get("/secret", func(c *fiber.Ctx) error {
		return c.SendString("love")
	})

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
		c.UUIDFactory,
	)

	http.NewAccount(accountSvc).Bind(app)
	http.NewRoom(roomSvc).Bind(app)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", config.C.Port)))
}
