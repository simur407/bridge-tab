package routes

import (
	"bridge-tab/api/middleware"
	tournament_management "bridge-tab/internal/tournament-management/application"
	domain "bridge-tab/internal/tournament-management/domain"

	"github.com/gofiber/fiber/v2"
)

func leaveTournament(repository domain.TournamentRepository) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		tournamentId := c.Params("tournamentId")
		contestantId := c.Locals("user").(middleware.UserMetadata).Id

		cmd := tournament_management.LeaveTournamentCommand{TournamentId: tournamentId, ContestantId: contestantId}
		err := cmd.Execute(repository)

		if err != nil {
			// TODO: handle error HTTP way
			return err
		}

		c.SendStatus(fiber.StatusOK)
		return nil
	}

}
