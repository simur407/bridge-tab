package middleware

import (
	"bridge-tab/internal/auth"

	"github.com/gofiber/fiber/v2"
)

type Token struct {
	Value string `cookie:"token"`
}

type UserMetadata struct {
	Id string
}

func JwtGuard() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := new(Token)
		if err := c.CookieParser(token); err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		id, err := auth.Decode(token.Value)
		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		c.Locals("user", UserMetadata{Id: id})

		return c.Next()
	}
}
