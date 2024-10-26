package middleware

import (
	"github.com/gofiber/fiber/v2"

	constants "github.com/flabio/safe_constants"
)

func ValidateToken(c *fiber.Ctx) error {
	token := c.Get(constants.AUTHORIZATION)
	if len(token) > 7 && token[:7] == constants.BEARER {
		return c.Next()
	}
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		constants.STATUS:  fiber.StatusUnauthorized,
		constants.MESSAGE: constants.TOKEN_INVALID,
	})
}
