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

type setRow struct {
	Id     string
	Label  string
	Boards string
}

func getSets(c *fiber.Ctx) error {
	s := newStores(c)
	tournament, err := loadTournament(c, s)
	if err != nil {
		return err
	}
	data, err := setsPage(s, tournament)
	if err != nil {
		return err
	}
	return view(c, "sets", data)
}

func postSet(c *fiber.Ctx) error {
	s := newStores(c)
	tournament, err := loadTournament(c, s)
	if err != nil {
		return err
	}
	data, err := setsPage(s, tournament)
	if err != nil {
		return err
	}
	label := strings.TrimSpace(c.FormValue("label"))
	data["FormLabel"] = label
	data["FormBoards"] = c.FormValue("boards")
	if label == "" {
		return reject(c, "sets", data, errors.New("etykieta jest pusta"))
	}
	boards, err := parseBoardNos(c.FormValue("boards"))
	if err != nil {
		return reject(c, "sets", data, err)
	}
	cmd := tournament_cmd.CreateSetCommand{
		TournamentId: tournament.Id,
		SetId:        uuid.NewString(),
		Label:        label,
		BoardNos:     boards,
	}
	if err := cmd.Execute(s.tournament); err != nil {
		return reject(c, "sets", data, err)
	}
	return c.Redirect("/tournaments/"+tournament.Id+"/sets", fiber.StatusSeeOther)
}

func postRemoveSet(c *fiber.Ctx) error {
	s := newStores(c)
	tournament, err := loadTournament(c, s)
	if err != nil {
		return err
	}
	data, err := setsPage(s, tournament)
	if err != nil {
		return err
	}
	cmd := tournament_cmd.RemoveSetCommand{TournamentId: tournament.Id, SetId: c.Params("setId")}
	if err := cmd.Execute(s.tournament); err != nil {
		return reject(c, "sets", data, err)
	}
	return c.Redirect("/tournaments/"+tournament.Id+"/sets", fiber.StatusSeeOther)
}

func setsPage(s stores, tournament *tournament_domain.TournamentDto) (fiber.Map, error) {
	sets, err := (&tournament_query.ListSetsQuery{TournamentId: tournament.Id}).Execute(s.sets)
	if err != nil {
		return nil, err
	}
	rows := make([]setRow, len(sets))
	for i, set := range sets {
		rows[i] = setRow{Id: set.Id, Label: set.Label, Boards: joinInts(set.BoardNos)}
	}
	data := pageData(tournament)
	data["Sets"] = rows
	data["FormLabel"] = ""
	data["FormBoards"] = ""
	return data, nil
}
