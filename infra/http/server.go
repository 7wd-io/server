package http

import (
	"7wd.io/rr"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"log/slog"
	"net/http"
)

func NewApp() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "7wd",
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			slog.Error(err.Error())

			var fe *fiber.Error
			if errors.As(err, &fe) {
				return ctx.Status(fe.Code).JSON(map[string]string{
					"errMessage": fe.Error(),
				})
			}

			var er rr.AppError
			if errors.As(err, &er) {
				return ctx.Status(http.StatusBadRequest).JSON(er)
			}

			return ctx.Status(http.StatusInternalServerError).SendString("internal server error")
		},
	})

	app.Use(logger.New())
	app.Use(recover.New())

	app.Use(useJWT())
	app.Use(useCORS())

	return app
}
