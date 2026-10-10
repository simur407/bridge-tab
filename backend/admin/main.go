package main

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"runtime"

	rounds_infra "bridge-tab/internal/rounds-registration/infrastructure"
	score_infra "bridge-tab/internal/score-registration/infrastructure"
	tournament_infra "bridge-tab/internal/tournament-management/infrastructure"
	user_infra "bridge-tab/internal/user/infrastructure"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	_ "github.com/lib/pq"
)

var (
	adminPassword string
	cookieSecure  bool
)

func main() {
	adminPassword = os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		panic("ADMIN_PASSWORD is required")
	}
	cookieSecure = os.Getenv("ADMIN_COOKIE_SECURE") == "1"

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("could not get current filename")
	}
	frontendDirectory := filepath.Join(filepath.Dir(filename), "frontend")

	app := fiber.New(fiber.Config{
		Views: html.New(frontendDirectory, ".html"),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var fiberErr *fiber.Error
			if errors.As(err, &fiberErr) {
				code = fiberErr.Code
			}
			return c.Status(code).SendString(err.Error())
		},
	})
	app.Use(logger.New())
	app.Use(recover.New())
	app.Static("/css", filepath.Join(frontendDirectory, "css"))

	dbString := os.Getenv("DATABASE_STRING")
	db, err := sql.Open("postgres", dbString)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		log.Debug(err)
		panic("failed to connect to database")
	}

	tournament_infra.Migrate(db)
	user_infra.Migrate(db)
	rounds_infra.Migrate(db)
	score_infra.Migrate(db)

	app.Get("/login", getLogin)
	app.Post("/login", postLogin)

	authed := app.Group("", requireAdmin, readTx(db), writeTx(db))
	authed.Post("/logout", postLogout)
	authed.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/tournaments", fiber.StatusSeeOther)
	})
	authed.Get("/tournaments", getTournaments)
	authed.Post("/tournaments", postTournament)
	authed.Get("/tournaments/:id", getTournament)
	authed.Post("/tournaments/:id/start", postStart)
	authed.Post("/tournaments/:id/finish", postFinish)

	authed.Get("/tournaments/:id/teams", getTeams)
	authed.Post("/tournaments/:id/teams", postTeam)
	authed.Post("/tournaments/:id/teams/:teamId/rename", postRenameTeam)
	authed.Post("/tournaments/:id/teams/:teamId/remove", postRemoveTeam)

	authed.Get("/tournaments/:id/tables", getTables)
	authed.Post("/tournaments/:id/tables", postTable)
	authed.Post("/tournaments/:id/tables/:tableId/remove", postRemoveTable)

	authed.Get("/tournaments/:id/protocols", getProtocols)
	authed.Post("/tournaments/:id/protocols", postProtocol)
	authed.Post("/tournaments/:id/protocols/:board/remove", postRemoveProtocol)

	authed.Get("/tournaments/:id/sets", getSets)
	authed.Post("/tournaments/:id/sets", postSet)
	authed.Post("/tournaments/:id/sets/:setId/remove", postRemoveSet)

	authed.Get("/tournaments/:id/rounds", getRounds)
	authed.Get("/tournaments/:id/rounds/:deal/:ns/:ew", getRound)
	authed.Post("/tournaments/:id/rounds/:deal/:ns/:ew", postRound)

	authed.Get("/tournaments/:id/scores", getScores)
	authed.Post("/tournaments/:id/scores/calculate", postCalculate)
	authed.Post("/tournaments/:id/scores/finish", postScoresFinish)

	port := os.Getenv("ADMIN_PORT")
	if port == "" {
		port = "3001"
	}
	if err := app.Listen(":" + port); err != nil {
		panic(err)
	}
}
