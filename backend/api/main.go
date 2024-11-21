package main

import (
	routes "bridge-tab/api/routes/tournament-management"

	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/lib/pq"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("I am alive")
	})

	dbString := "postgres://bridge-tab:bridge-tab@localhost/bridge-tab?sslmode=disable"
	db, err := sql.Open("postgres", dbString)

	if err != nil {
		panic(err)
	}

	routes.TournamentManagement(app, db)

	app.Listen(":3000")
}
