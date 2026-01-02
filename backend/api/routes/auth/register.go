package routes

import (
	"bridge-tab/api/middleware"
	application "bridge-tab/internal/user/application"
	infra "bridge-tab/internal/user/infrastructure"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type RegisterUserRequestDto struct {
	Name string `json:"name"`
}

func register() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		body := new(RegisterUserRequestDto)

		if err := c.BodyParser(body); err != nil {
			return err
		}

		tx := middleware.GetTransaction(c)
		repository := infra.PostgresUserRepository{Ctx: c.UserContext(), Tx: tx}

		id := uuid.New().String()

		command := &application.RegisterUserCommand{
			Id:   id,
			Name: body.Name,
		}
		if err := command.Execute(&repository); err != nil {
			return err
		}

		return registerResponse(c, id)
	}
}

func registerResponse(c *fiber.Ctx, id string) error {
	response := fiber.Map{
		"login": id,
	}
	return c.JSON(response)
}
