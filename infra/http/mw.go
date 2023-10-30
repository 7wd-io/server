package http

import (
	"7wd.io/config"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"net/http"
	"strings"
)

func useCORS() fiber.Handler {
	methods := strings.Join([]string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPut,
		http.MethodPatch,
		http.MethodPost,
		http.MethodDelete,
		http.MethodOptions,
	}, ",")

	return cors.New(cors.Config{
		AllowOrigins:     config.C.ClientOrigin,
		AllowMethods:     methods,
		AllowCredentials: true,
	})
}

func useJWT() fiber.Handler {
	return jwtware.New(jwtware.Config{
		Filter: func(ctx *fiber.Ctx) bool {
			switch ctx.Path() {
			case
				"/ping",
				"/account/signup",
				"/account/signin",
				"/account/refresh":
				return true
			}

			return false
		},
		SigningKey:   jwtware.SigningKey{Key: []byte(config.C.Secret)},
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(errUnauthorized)
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}
