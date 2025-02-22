package http

import (
	"7wd.io/domain"
	"7wd.io/rr"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"log/slog"
	"net/http"
)

func New() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "7wd",
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			slog.Error(err.Error())

			var fe *fiber.Error
			if errors.As(err, &fe) {
				return ctx.Status(fe.Code).JSON(map[string]string{
					"err": fe.Error(),
				})
			}

			var er rr.AppError
			if errors.As(err, &er) {
				if errors.Is(er, errUnauthorized) || errors.Is(er, domain.ErrSessionNotFound) {
					return ctx.Status(http.StatusUnauthorized).JSON(er)
				}

				return ctx.Status(http.StatusBadRequest).JSON(er)
			}

			return ctx.Status(http.StatusInternalServerError).JSON(map[string]string{
				"err": "internal server error",
			})
		},
	})

	app.Use(logger.New())
	app.Use(recover.New())

	app.Use(useJWT())
	//app.Use(useCORS())

	return app
}
