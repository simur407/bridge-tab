package middleware

import (
	"bridge-tab/internal/auth"
	"bridge-tab/internal/basepath"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// BasePath is the optional public URL prefix for this HTTP frontend.
var BasePath string

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
			return c.Redirect(registerURL(c.Path()))
		}

		id, err := auth.Decode(token.Value)
		if err != nil {
			log.Debug(err)
			clearCookie(c)
			return c.Redirect(registerURL(c.Path()))
		}

		c.Locals("user", UserMetadata{Id: id})

		return c.Next()
	}
}

func registerURL(redirectPath string) string {
	return basepath.Join(BasePath, "/register") + "?redirect=" + url.QueryEscape(redirectPath)
}

func clearCookie(c *fiber.Ctx) {
	cookie := new(fiber.Cookie)
	cookie.Name = "token"
	cookie.Expires = time.Now()
	cookie.Path = basepath.CookiePath(BasePath)
	c.Cookie(cookie)
}
