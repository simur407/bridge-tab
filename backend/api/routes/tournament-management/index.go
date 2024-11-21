package routes

import (
	"database/sql"

	infra "bridge-tab/internal/tournament-management/infrastructure"

	"github.com/gofiber/fiber/v2"
)

func TournamentManagement(app fiber.Router, db *sql.DB) error {
	index := app.Group("/tournament-management")

	repository := infra.PostgresTournamentRepository{Db: db}
	index.Get("/tournaments", )
	index.Post("/tournaments/:tournamentId/join", joinTournament(&repository))
	index.Post("/tournaments/:tournamentId/leave", leaveTournament(&repository))
	index.Post("/tournaments/:tournamentId/teams/:teamId/join", joinTeam(&repository))
	index.Post("/tournaments/:tournamentId/teams/:teamId/leave", leaveTeam(&repository))
	return nil
}
