package http

import (
	"7wd.io/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func usePassport(c *fiber.Ctx) (p domain.Passport, err error) {
	t, ok := c.Locals("user").(*jwt.Token)

	if !ok {
		return p, errInvalidToken
	}

	claims, ok := t.Claims.(jwt.MapClaims)

	if !ok {
		return p, errInvalidToken
	}

	return domain.Passport{
		Id:       domain.UserId(claims["id"].(float64)),
		Nickname: domain.Nickname(claims["nickname"].(string)),
		Rating:   claims["rating"].(domain.Rating),
	}, nil
}
