package main

import (
	"strings"

	tournament_cmd "bridge-tab/internal/tournament-management/application/command"
	tournament_query "bridge-tab/internal/tournament-management/application/query"
	tournament_domain "bridge-tab/internal/tournament-management/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type tableRow struct {
	Id     string
	Number int
}

func getTables(c *fiber.Ctx) error {
	s := newStores(c)
	tournament, err := loadTournament(c, s)
	if err != nil {
		return err
	}
	data, err := tablesPage(s, tournament)
	if err != nil {
		return err
	}
	return view(c, "tables", data)
}

func postTable(c *fiber.Ctx) error {
	s := newStores(c)
	tournament, err := loadTournament(c, s)
	if err != nil {
		return err
	}
	data, err := tablesPage(s, tournament)
	if err != nil {
		return err
	}
	number, err := optionalPositive(c.FormValue("number"))
	data["FormNumber"] = strings.TrimSpace(c.FormValue("number"))
	if err != nil {
		return reject(c, "tables", data, err)
	}
	cmd := tournament_cmd.CreateTableCommand{
		TournamentId: tournament.Id,
		TableId:      uuid.NewString(),
		Number:       number,
	}
	if err := cmd.Execute(s.tournament); err != nil {
		return reject(c, "tables", data, err)
	}
	return c.Redirect("/tournaments/"+tournament.Id+"/tables", fiber.StatusSeeOther)
}

func postRemoveTable(c *fiber.Ctx) error {
	s := newStores(c)
	tournament, err := loadTournament(c, s)
	if err != nil {
		return err
	}
	data, err := tablesPage(s, tournament)
	if err != nil {
		return err
	}
	cmd := tournament_cmd.RemoveTableCommand{TournamentId: tournament.Id, TableId: c.Params("tableId")}
	if err := cmd.Execute(s.tournament); err != nil {
		return reject(c, "tables", data, err)
	}
	return c.Redirect("/tournaments/"+tournament.Id+"/tables", fiber.StatusSeeOther)
}

func tablesPage(s stores, tournament *tournament_domain.TournamentDto) (fiber.Map, error) {
	tables, err := (&tournament_query.ListTablesQuery{TournamentId: tournament.Id}).Execute(s.tables)
	if err != nil {
		return nil, err
	}
	rows := make([]tableRow, len(tables))
	for i, table := range tables {
		rows[i] = tableRow{Id: table.Id, Number: table.Number}
	}
	data := pageData(tournament)
	data["Tables"] = rows
	data["FormNumber"] = ""
	return data, nil
}
