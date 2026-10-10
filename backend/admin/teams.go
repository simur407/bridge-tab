package main

import (
	"errors"
	"strings"

	tournament_cmd "bridge-tab/internal/tournament-management/application/command"
	tournament_query "bridge-tab/internal/tournament-management/application/query"
	tournament_domain "bridge-tab/internal/tournament-management/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type teamRow struct {
	Id      string
	Number  int
	Name    string
	Players string
}

func getTeams(c *fiber.Ctx) error {
	s := newStores(c)
	tournament, err := loadTournament(c, s)
	if err != nil {
		return err
	}
	data, err := teamsPage(s, tournament)
	if err != nil {
		return err
	}
	return view(c, "teams", data)
}

func postTeam(c *fiber.Ctx) error {
	s := newStores(c)
	tournament, err := loadTournament(c, s)
	if err != nil {
		return err
	}
	data, err := teamsPage(s, tournament)
	if err != nil {
		return err
	}
	name := strings.TrimSpace(c.FormValue("name"))
	number, err := optionalPositive(c.FormValue("number"))
	data["FormName"] = name
	data["FormNumber"] = strings.TrimSpace(c.FormValue("number"))
	if err != nil {
		return reject(c, "teams", data, err)
	}
	if name == "" {
		return reject(c, "teams", data, errors.New("nazwa jest pusta"))
	}
	teams, err := (&tournament_query.ListTeamsQuery{TournamentId: tournament.Id}).Execute(s.teams)
	if err != nil {
		return err
	}
	for _, team := range teams {
		if team.Name == name {
			return reject(c, "teams", data, tournament_domain.ErrTeamNameAlreadyExists)
		}
	}
	cmd := tournament_cmd.CreateTeamCommand{
		TournamentId: tournament.Id,
		TeamId:       uuid.NewString(),
		Name:         name,
		Number:       number,
	}
	if err := cmd.Execute(s.tournament); err != nil {
		return reject(c, "teams", data, err)
	}
	return c.Redirect("/tournaments/"+tournament.Id+"/teams", fiber.StatusSeeOther)
}

func postRenameTeam(c *fiber.Ctx) error {
	s := newStores(c)
	tournament, err := loadTournament(c, s)
	if err != nil {
		return err
	}
	data, err := teamsPage(s, tournament)
	if err != nil {
		return err
	}
	cmd := tournament_cmd.RenameTeamCommand{
		TournamentId: tournament.Id,
		TeamId:       c.Params("teamId"),
		Name:         c.FormValue("name"),
	}
	if err := cmd.Execute(s.tournament); err != nil {
		return reject(c, "teams", data, err)
	}
	return c.Redirect("/tournaments/"+tournament.Id+"/teams", fiber.StatusSeeOther)
}

func postRemoveTeam(c *fiber.Ctx) error {
	s := newStores(c)
	tournament, err := loadTournament(c, s)
	if err != nil {
		return err
	}
	data, err := teamsPage(s, tournament)
	if err != nil {
		return err
	}
	cmd := tournament_cmd.RemoveTeamCommand{TournamentId: tournament.Id, TeamId: c.Params("teamId")}
	if err := cmd.Execute(s.tournament); err != nil {
		return reject(c, "teams", data, err)
	}
	return c.Redirect("/tournaments/"+tournament.Id+"/teams", fiber.StatusSeeOther)
}

func teamsPage(s stores, tournament *tournament_domain.TournamentDto) (fiber.Map, error) {
	teams, err := (&tournament_query.ListTeamsQuery{TournamentId: tournament.Id}).Execute(s.teams)
	if err != nil {
		return nil, err
	}
	names, err := nameMap(s)
	if err != nil {
		return nil, err
	}
	rows := make([]teamRow, len(teams))
	for i, team := range teams {
		rows[i] = teamRow{
			Id:      team.Id,
			Number:  team.Number,
			Name:    team.Name,
			Players: playerNames(team.Members, names),
		}
	}
	data := pageData(tournament)
	data["Teams"] = rows
	data["FormName"] = ""
	data["FormNumber"] = ""
	return data, nil
}
