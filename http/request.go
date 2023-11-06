package http

import (
	"7wd.io/adapter/validator"
	"github.com/gofiber/fiber/v2"
)

func useQueryRequest(c *fiber.Ctx, req interface{}) error {
	if err := c.QueryParser(req); err != nil {
		return err
	}

	return validator.Validate(req)
}

func useBodyRequest(c *fiber.Ctx, req interface{}) error {
	if err := c.BodyParser(req); err != nil {
		return err
	}

	if err := validator.Validate(req); err != nil {
		return errInvalidRequest
	}

	return nil
}
