package routes

import (
	"time"

	"bridge-tab/api/middleware"
	auth "bridge-tab/internal/auth"
	infra "bridge-tab/internal/user/infrastructure"

	"github.com/gofiber/fiber/v2"
)

type LoginRequestDto struct {
	Login string `json:"login"`
	// password will be done in the future
}

func login() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		body := new(LoginRequestDto)

		if err := c.BodyParser(body); err != nil {
			return err
		}

		tx := middleware.GetTransaction(c)
		repository := infra.PostgresUserRepository{Tx: tx}
		_, err := repository.GetById(body.Login)
		if err != nil {
			return err
		}

		token, err := auth.Generate(body.Login)

		if err != nil {
			return err
		}

		cookie := new(fiber.Cookie)
		cookie.Name = "token"
		cookie.Value = token
		cookie.Expires = time.Now().Add(auth.EXPIRES_AT * 10)
		c.Cookie(cookie)

		c.SendStatus(fiber.StatusOK)
		return nil
	}
}
