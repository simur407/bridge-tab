package routes

import (
	"bridge-tab/api/middleware"
	application "bridge-tab/internal/tournament-management/application"
	domain "bridge-tab/internal/tournament-management/domain"
	infra "bridge-tab/internal/tournament-management/infrastructure"

	"github.com/gofiber/fiber/v2"
)

func findTeamByName() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		tournamentId := c.Params("tournamentId")
		teamId := c.Queries()["name"]

		tx := middleware.GetTransaction(c)
		repository := infra.PostgresTeamReadRepository{Ctx: c.UserContext(), Tx: tx}

		cmd := application.GetTeamByNameQuery{
			TournamentId: tournamentId,
			Name:         teamId,
		}
		t, err := cmd.Execute(&repository)

		if err != nil {
			// TODO: handle error HTTP way
			return err
		}

		return findTeamByNameResponse(c, *t)
	}
}

func findTeamByNameResponse(c *fiber.Ctx, team domain.TeamDto) error {
	response := fiber.Map{
		"id":   team.Id,
		"name": team.Name,
	}
	return c.JSON(response)
}
