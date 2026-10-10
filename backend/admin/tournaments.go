package main

import (
	"errors"
	"sort"
	"strings"

	tournament_cmd "bridge-tab/internal/tournament-management/application/command"
	tournament_query "bridge-tab/internal/tournament-management/application/query"
	tournament_domain "bridge-tab/internal/tournament-management/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type tournamentRow struct {
	Id     string
	Name   string
	Status string
	Code   string
}

type contestantRow struct {
	Name string
	Id   string
}

func getTournaments(c *fiber.Ctx) error {
	s := newStores(c)
	list, err := (&tournament_query.ListTournamentsQuery{}).Execute(s.tournamentRead)
	if err != nil {
		return err
	}
	return view(c, "tournaments", fiber.Map{
		"Title":       "Turnieje",
		"Tournaments": tournamentRows(list),
		"FormName":    "",
	})
}

func postTournament(c *fiber.Ctx) error {
	s := newStores(c)
	list, err := (&tournament_query.ListTournamentsQuery{}).Execute(s.tournamentRead)
	if err != nil {
		return err
	}
	data := fiber.Map{
		"Title":       "Turnieje",
		"Tournaments": tournamentRows(list),
		"FormName":    strings.TrimSpace(c.FormValue("name")),
	}
	name := data["FormName"].(string)
	if name == "" {
		return reject(c, "tournaments", data, errors.New("nazwa jest pusta"))
	}
	cmd := tournament_cmd.CreateTournamentCommand{TournamentId: uuid.NewString(), Name: name}
	if err := cmd.Execute(s.tournament); err != nil {
		return reject(c, "tournaments", data, err)
	}
	return c.Redirect("/tournaments/"+cmd.TournamentId, fiber.StatusSeeOther)
}

func tournamentRows(list []tournament_domain.TournamentDto) []tournamentRow {
	rows := make([]tournamentRow, len(list))
	for i, tournament := range list {
		rows[i] = tournamentRow{
			Id:     tournament.Id,
			Name:   tournament.Name,
			Status: statusLabel(tournament),
			Code:   tournament.Id,
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Name < rows[j].Name
	})
	return rows
}

func getTournament(c *fiber.Ctx) error {
	s := newStores(c)
	tournament, err := loadTournament(c, s)
	if err != nil {
		return err
	}
	data, err := tournamentPage(s, tournament)
	if err != nil {
		return err
	}
	return view(c, "tournament", data)
}

func postStart(c *fiber.Ctx) error {
	s := newStores(c)
	tournament, err := loadTournament(c, s)
	if err != nil {
		return err
	}
	data, err := tournamentPage(s, tournament)
	if err != nil {
		return err
	}
	cmd := tournament_cmd.StartTurnamentCommand{TournamentId: tournament.Id}
	if err := cmd.Execute(s.tournament, s.teams, s.boards, s.sessions); err != nil {
		return reject(c, "tournament", data, err)
	}
	return c.Redirect("/tournaments/"+tournament.Id, fiber.StatusSeeOther)
}

func postFinish(c *fiber.Ctx) error {
	return finishTournament(c, "/tournaments/"+c.Params("id"))
}

func finishTournament(c *fiber.Ctx, back string) error {
	s := newStores(c)
	tournament, err := loadTournament(c, s)
	if err != nil {
		return err
	}
	cmd := tournament_cmd.FinishTournamentCommand{TournamentId: tournament.Id}
	if err := cmd.Execute(s.tournament, s.sessions, s.boards, s.played); err != nil {
		if strings.HasSuffix(back, "/scores") {
			data, pageErr := scoresPage(s, tournament)
			if pageErr != nil {
				return pageErr
			}
			return reject(c, "scores", data, err)
		}
		data, pageErr := tournamentPage(s, tournament)
		if pageErr != nil {
			return pageErr
		}
		return reject(c, "tournament", data, err)
	}
	return c.Redirect(back, fiber.StatusSeeOther)
}

func tournamentPage(s stores, tournament *tournament_domain.TournamentDto) (fiber.Map, error) {
	contestants, err := (&tournament_query.ListContestantsQuery{TournamentId: tournament.Id}).Execute(s.tournamentRead)
	if err != nil {
		return nil, err
	}
	names, err := nameMap(s)
	if err != nil {
		return nil, err
	}
	rows := make([]contestantRow, len(contestants))
	for i, contestant := range contestants {
		id := string(contestant.Id)
		name := names[strings.ToLower(id)]
		if name == "" {
			name = id
		}
		rows[i] = contestantRow{Name: name, Id: id}
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Name < rows[j].Name
	})
	data := pageData(tournament)
	data["Contestants"] = rows
	return data, nil
}
