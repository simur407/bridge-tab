package main

import (
	"fmt"
	"strconv"
	"strings"

	tournament_cmd "bridge-tab/internal/tournament-management/application/command"
	tournament_query "bridge-tab/internal/tournament-management/application/query"
	tournament_domain "bridge-tab/internal/tournament-management/domain"

	"github.com/gofiber/fiber/v2"
)

type scheduleRow struct {
	RoundNo     int
	TableNumber int
	NsNumber    int
	EwNumber    int
	SetLabel    string
	TableId     string
}

func getSchedule(c *fiber.Ctx) error {
	s := newStores(c)
	tournament, err := loadTournament(c, s)
	if err != nil {
		return err
	}
	data, err := schedulePage(s, tournament)
	if err != nil {
		return err
	}
	return view(c, "schedule", data)
}

func postSchedule(c *fiber.Ctx) error {
	s := newStores(c)
	tournament, err := loadTournament(c, s)
	if err != nil {
		return err
	}
	data, err := schedulePage(s, tournament)
	if err != nil {
		return err
	}
	data["FormRound"] = strings.TrimSpace(c.FormValue("round"))
	data["FormTable"] = strings.TrimSpace(c.FormValue("table"))
	data["FormNS"] = strings.TrimSpace(c.FormValue("ns"))
	data["FormEW"] = strings.TrimSpace(c.FormValue("ew"))
	data["FormSet"] = strings.TrimSpace(c.FormValue("set"))

	roundNo, err := strconv.Atoi(strings.TrimSpace(c.FormValue("round")))
	if err != nil || roundNo < 1 {
		return reject(c, "schedule", data, fmt.Errorf("zły numer rundy"))
	}
	tableNumber, err := strconv.Atoi(strings.TrimSpace(c.FormValue("table")))
	if err != nil || tableNumber < 1 {
		return reject(c, "schedule", data, fmt.Errorf("zły numer stołu"))
	}
	nsNumber, err := strconv.Atoi(strings.TrimSpace(c.FormValue("ns")))
	if err != nil || nsNumber < 1 {
		return reject(c, "schedule", data, fmt.Errorf("zły numer NS"))
	}
	ewNumber, err := strconv.Atoi(strings.TrimSpace(c.FormValue("ew")))
	if err != nil || ewNumber < 1 {
		return reject(c, "schedule", data, fmt.Errorf("zły numer EW"))
	}
	setLabel := strings.TrimSpace(c.FormValue("set"))
	if setLabel == "" {
		return reject(c, "schedule", data, fmt.Errorf("etykieta seta jest pusta"))
	}

	table, err := (&tournament_query.GetTableByNumberQuery{TournamentId: tournament.Id, Number: tableNumber}).Execute(s.tables)
	if err != nil {
		return reject(c, "schedule", data, fmt.Errorf("stół %d: %w", tableNumber, explain(err)))
	}
	ns, err := (&tournament_query.GetTeamByNumberQuery{TournamentId: tournament.Id, Number: nsNumber}).Execute(s.teams)
	if err != nil {
		return reject(c, "schedule", data, fmt.Errorf("drużyna NS %d: %w", nsNumber, explain(err)))
	}
	ew, err := (&tournament_query.GetTeamByNumberQuery{TournamentId: tournament.Id, Number: ewNumber}).Execute(s.teams)
	if err != nil {
		return reject(c, "schedule", data, fmt.Errorf("drużyna EW %d: %w", ewNumber, explain(err)))
	}
	sets, err := (&tournament_query.ListSetsQuery{TournamentId: tournament.Id}).Execute(s.sets)
	if err != nil {
		return err
	}
	var setId string
	for _, set := range sets {
		if strings.EqualFold(set.Label, setLabel) {
			setId = set.Id
			break
		}
	}
	if setId == "" {
		return reject(c, "schedule", data, fmt.Errorf("nie znaleziono seta %s", setLabel))
	}

	cmd := tournament_cmd.AddScheduleCommand{
		TournamentId: tournament.Id,
		RoundNo:      roundNo,
		TableId:      table.Id,
		NsTeamId:     ns.Id,
		EwTeamId:     ew.Id,
		SetId:        setId,
	}
	if err := cmd.Execute(s.tournament); err != nil {
		return reject(c, "schedule", data, err)
	}
	return c.Redirect("/tournaments/"+tournament.Id+"/schedule", fiber.StatusSeeOther)
}

func postRemoveSchedule(c *fiber.Ctx) error {
	s := newStores(c)
	tournament, err := loadTournament(c, s)
	if err != nil {
		return err
	}
	data, err := schedulePage(s, tournament)
	if err != nil {
		return err
	}
	roundNo, err := pathInt(c, "round")
	if err != nil {
		return err
	}
	cmd := tournament_cmd.RemoveScheduleCommand{
		TournamentId: tournament.Id,
		RoundNo:      roundNo,
		TableId:      c.Params("tableId"),
	}
	if err := cmd.Execute(s.tournament); err != nil {
		return reject(c, "schedule", data, err)
	}
	return c.Redirect("/tournaments/"+tournament.Id+"/schedule", fiber.StatusSeeOther)
}

func schedulePage(s stores, tournament *tournament_domain.TournamentDto) (fiber.Map, error) {
	entries, err := (&tournament_query.ListScheduleQuery{TournamentId: tournament.Id}).Execute(s.schedule)
	if err != nil {
		return nil, err
	}
	rows := make([]scheduleRow, len(entries))
	for i, entry := range entries {
		rows[i] = scheduleRow{
			RoundNo:     entry.RoundNo,
			TableNumber: entry.TableNumber,
			NsNumber:    entry.NsNumber,
			EwNumber:    entry.EwNumber,
			SetLabel:    entry.SetLabel,
			TableId:     entry.TableId,
		}
	}
	data := pageData(tournament)
	data["Schedule"] = rows
	data["FormRound"] = ""
	data["FormTable"] = ""
	data["FormNS"] = ""
	data["FormEW"] = ""
	data["FormSet"] = ""
	return data, nil
}
