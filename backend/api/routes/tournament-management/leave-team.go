package routes

import (
	tournament_management "bridge-tab/internal/tournament-management/application"
	domain "bridge-tab/internal/tournament-management/domain"

	"github.com/gofiber/fiber/v2"
)

type LeaveTeamRequestDto struct {
	ContestantId string `json:"contestantId"`
}

func leaveTeam(repository domain.TournamentRepository) func (c *fiber.Ctx) error {
	return func (c *fiber.Ctx) error {
		tournamentId := c.Params("tournamentId")
		teamId := c.Params("teamId")

		body := new(LeaveTeamRequestDto)

		if err := c.BodyParser(body); err != nil {
			return err
		}

		cmd := tournament_management.LeaveTeamCommand{ 
			TournamentId: tournamentId,
			TeamId: teamId,
			ContestantId: body.ContestantId,
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
