package main

import (
	auth "bridge-tab/api/routes/auth"
	tournament_management "bridge-tab/api/routes/tournament-management"

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

	app.Static("/", "../frontend")

	dbString := "postgres://bridge-tab:bridge-tab@localhost/bridge-tab?sslmode=disable"
	db, err := sql.Open("postgres", dbString)

	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic("failed to connect to database")
	}

	auth.Auth(app, db)
	tournament_management.TournamentManagement(app, db)

	app.Listen(":3000")
}
