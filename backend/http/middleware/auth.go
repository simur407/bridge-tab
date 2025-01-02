package middleware

import (
	"bridge-tab/internal/auth"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
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
			log.Debug(err)
			clearCookie(c)
			return c.Redirect("/register?redirect=" + c.Path())
		}

		id, err := auth.Decode(token.Value)
		if err != nil {
			log.Debug(err)
			clearCookie(c)
			return c.Redirect("/register?redirect=" + c.Path())
		}

		c.Locals("user", UserMetadata{Id: id})

		return c.Next()
	}
}

func clearCookie(c *fiber.Ctx) {
	cookie := new(fiber.Cookie)
	cookie.Name = "token"
	cookie.Expires = time.Now()
	c.Cookie(cookie)
}
