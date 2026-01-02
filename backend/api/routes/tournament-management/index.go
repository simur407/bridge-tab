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

	// index.Get("/tournaments")
	index.Post("/tournaments/:tournamentId/join", middleware.Transaction(db, nil), joinTournament())
	index.Post("/tournaments/:tournamentId/leave", middleware.Transaction(db, nil), leaveTournament())
	index.Get("/tournaments/:tournamentId/teams", middleware.Transaction(db, &sql.TxOptions{ReadOnly: true}), findTeamByName())
	index.Post("/tournaments/:tournamentId/teams/:teamId/join", middleware.Transaction(db, nil), joinTeam())
	index.Post("/tournaments/:tournamentId/teams/:teamId/leave", middleware.Transaction(db, nil), leaveTeam())
	return nil
}
