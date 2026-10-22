package main

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"runtime"

	play_infra "bridge-tab/internal/play/infrastructure"
	score_infra "bridge-tab/internal/score-registration/infrastructure"
	tournament_infra "bridge-tab/internal/tournament-management/infrastructure"
	umpirecall "bridge-tab/internal/umpire-call"
	user_infra "bridge-tab/internal/user/infrastructure"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	_ "github.com/lib/pq"
)

var (
	tabletPassword string
	cookieSecure   bool
)

func main() {
	tabletPassword = os.Getenv("TABLET_PASSWORD")
	if tabletPassword == "" {
		panic("TABLET_PASSWORD is required")
	}
	cookieSecure = os.Getenv("TABLET_COOKIE_SECURE") == "1"

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
	app.Static("/script", filepath.Join(frontendDirectory, "script"))

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
	play_infra.Migrate(db)
	score_infra.Migrate(db)
	umpirecall.Migrate(db)

	app.Get("/login", getLogin)
	app.Post("/login", postLogin)

	authed := app.Group("", requireTablet, readTx(db), writeTx(db))
	authed.Post("/logout", postLogout)
	authed.Get("/", getHome)
	authed.Get("/setup", getSetup)
	authed.Post("/setup", postSetup)
	authed.Get("/setup/tables", getSetupTables)
	authed.Get("/play", getPlay)
	authed.Get("/play/umpire-call", getUmpireCall)
	authed.Post("/play/umpire-call", postUmpireCall)
	authed.Post("/play/advance", postAdvance)
	authed.Get("/play/boards/:dealNo", getBoard)
	authed.Post("/play/verify", postVerify)
	authed.Post("/play/add", postAdd)
	authed.Get("/play/boards/:dealNo/results", getResults)

	port := os.Getenv("TABLET_PORT")
	if port == "" {
		port = "3002"
	}
	if err := app.Listen(":" + port); err != nil {
		panic(err)
	}
}
