package routes

import (
	"database/sql"

	infra "bridge-tab/internal/user/infrastructure"

	"github.com/gofiber/fiber/v2"
)

func Auth(app fiber.Router, db *sql.DB) error {
	index := app.Group("/auth")

	repository := infra.PostgresUserRepository{Db: db}
	index.Post("/login", login())
	index.Post("/register", register(&repository))
	return nil
}
