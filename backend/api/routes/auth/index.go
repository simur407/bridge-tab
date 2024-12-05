package routes

import (
	"database/sql"

	"bridge-tab/api/middleware"
	infra "bridge-tab/internal/user/infrastructure"

	"github.com/gofiber/fiber/v2"
)

func Auth(app fiber.Router, db *sql.DB) error {
	infra.Migrate(db)

	index := app.Group("/auth")

	index.Post("/login", middleware.Transaction(db, &sql.TxOptions{ReadOnly: true}), login())
	index.Post("/register", middleware.Transaction(db, nil), register())
	return nil
}
