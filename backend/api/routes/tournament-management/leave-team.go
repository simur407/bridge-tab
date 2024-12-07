package routes

import (
	"bridge-tab/api/middleware"
	application "bridge-tab/internal/tournament-management/application"
	domain "bridge-tab/internal/tournament-management/domain"

	"github.com/gofiber/fiber/v2"
)

type LeaveTeamRequestDto struct {
	ContestantId string `json:"contestantId"`
}

func leaveTeam(repository domain.TournamentRepository) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		tournamentId := c.Params("tournamentId")
		teamId := c.Params("teamId")
		contestantId := c.Locals("user").(middleware.UserMetadata).Id

		cmd := application.LeaveTeamCommand{
			TournamentId: tournamentId,
			TeamId:       teamId,
			ContestantId: contestantId,
		}
		err := cmd.Execute(repository)

		if err != nil {
			// TODO: handle error HTTP way
			return err
		}

		c.SendStatus(fiber.StatusOK)
		return nil
	}

}
