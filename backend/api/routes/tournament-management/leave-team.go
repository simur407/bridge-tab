package routes

import (
	"bridge-tab/api/middleware"
	application "bridge-tab/internal/tournament-management/application"
	infra "bridge-tab/internal/tournament-management/infrastructure"

	"github.com/gofiber/fiber/v2"
)

type LeaveTeamRequestDto struct {
	ContestantId string `json:"contestantId"`
}

func leaveTeam() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		tournamentId := c.Params("tournamentId")
		teamId := c.Params("teamId")
		contestantId := c.Locals("user").(middleware.UserMetadata).Id

		tx := middleware.GetTransaction(c)
		repository := infra.PostgresTournamentRepository{Ctx: c.UserContext(), Tx: tx}

		cmd := application.LeaveTeamCommand{
			TournamentId: tournamentId,
			TeamId:       teamId,
			ContestantId: contestantId,
		}
		err := cmd.Execute(&repository)

		if err != nil {
			// TODO: handle error HTTP way
			return err
		}

		c.SendStatus(fiber.StatusOK)
		return nil
	}

}
