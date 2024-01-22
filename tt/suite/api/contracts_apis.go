package api

import (
	"github.com/gofiber/fiber/v2"
)

type Binder interface {
	Bind(app *fiber.App)
}
