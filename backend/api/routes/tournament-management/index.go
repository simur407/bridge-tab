package routes

import (
	"database/sql"

	middleware "bridge-tab/api/middleware"
	infra "bridge-tab/internal/tournament-management/infrastructure"

	"github.com/gofiber/fiber/v2"
)

func TournamentManagement(app fiber.Router, db *sql.DB) error {
	infra.Migrate(db)

	index := app.Group("/tournament-management")
	index.Use(middleware.JwtGuard())

	repository := infra.PostgresTournamentRepository{Db: db}
	// index.Get("/tournaments")
	index.Post("/tournaments/:tournamentId/join", joinTournament(&repository))
	index.Post("/tournaments/:tournamentId/leave", leaveTournament(&repository))
	index.Post("/tournaments/:tournamentId/teams/:teamId/join", joinTeam(&repository))
	index.Post("/tournaments/:tournamentId/teams/:teamId/leave", leaveTeam(&repository))
	return nil
}
