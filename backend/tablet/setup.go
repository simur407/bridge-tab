package main

import (
	"database/sql"
	"errors"
	"sort"
	"strconv"
	"strings"

	tournament_query "bridge-tab/internal/tournament-management/application/query"
	tournament_domain "bridge-tab/internal/tournament-management/domain"

	"github.com/gofiber/fiber/v2"
)

func getHome(c *fiber.Ctx) error {
	a, ok := readAssignment(c)
	if !ok {
		return c.Redirect("/setup", fiber.StatusSeeOther)
	}
	s := newStores(c)
	tournament, err := findTournament(s, a.TournamentId)
	if err != nil {
		return err
	}
	if tournament == nil {
		return view(c, "setup", fiber.Map{
			"Title": "Konfiguracja",
			"Error": "Nie znaleziono turnieju. Wybierz turniej i stół ponownie.",
		})
	}
	if tournament.StartedAt == "" {
		return view(c, "waiting", waitingData(tournament, a))
	}
	if isHTMX(c) {
		c.Set("HX-Redirect", "/play")
		return c.SendStatus(fiber.StatusNoContent)
	}
	return c.Redirect("/play", fiber.StatusSeeOther)
}

func getSetup(c *fiber.Ctx) error {
	s := newStores(c)
	data, err := setupPage(c, s, "")
	if err != nil {
		return err
	}
	return view(c, "setup", data)
}

func postSetup(c *fiber.Ctx) error {
	s := newStores(c)
	existing, hasAssignment := readAssignment(c)
	if hasAssignment {
		current, err := findTournament(s, existing.TournamentId)
		if err != nil {
			return err
		}
		if current != nil && assignmentLocked(existing, current.StartedAt) {
			data, pageErr := setupPage(c, s, "")
			if pageErr != nil {
				return pageErr
			}
			return reject(c, "setup", data, errors.New("turniej już trwa — nie można zmienić stołu na tym urządzeniu"))
		}
	}

	tournamentId := strings.TrimSpace(c.FormValue("tournamentId"))
	tableRaw := strings.TrimSpace(c.FormValue("tableNumber"))
	data, err := setupPage(c, s, tournamentId)
	if err != nil {
		return err
	}
	if tournamentId == "" {
		return reject(c, "setup", data, errors.New("wybierz turniej"))
	}
	tableNumber, err := strconv.Atoi(tableRaw)
	if err != nil || tableNumber < 1 {
		return reject(c, "setup", data, errors.New("wybierz stół"))
	}
	tournament, err := findTournament(s, tournamentId)
	if err != nil {
		return reject(c, "setup", data, err)
	}
	if tournament == nil {
		return reject(c, "setup", data, errors.New("nie znaleziono turnieju"))
	}
	table, err := s.tables.FindByNumber(&tournamentId, tableNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return reject(c, "setup", data, errors.New("nie znaleziono stołu"))
		}
		return reject(c, "setup", data, err)
	}
	if table == nil || table.Number != tableNumber {
		return reject(c, "setup", data, errors.New("nie znaleziono stołu"))
	}
	setAssignment(c, assignment{TournamentId: tournamentId, TableNumber: tableNumber})
	return c.Redirect("/", fiber.StatusSeeOther)
}

func getSetupTables(c *fiber.Ctx) error {
	s := newStores(c)
	tournamentId := strings.TrimSpace(c.Query("tournamentId"))
	tables, err := loadTables(s, tournamentId)
	if err != nil {
		return err
	}
	return fragment(c, "setup-tables", fiber.Map{
		"Tables":      tables,
		"TableNumber": 0,
	})
}

func setupPage(c *fiber.Ctx, s stores, selectedTournamentId string) (fiber.Map, error) {
	list, err := (&tournament_query.ListTournamentsQuery{}).Execute(s.tournamentRead)
	if err != nil {
		return nil, err
	}
	a, hasAssignment := readAssignment(c)
	locked := false
	tableNumber := 0
	if selectedTournamentId == "" && hasAssignment {
		selectedTournamentId = a.TournamentId
		tableNumber = a.TableNumber
	}
	if hasAssignment {
		current, err := findTournament(s, a.TournamentId)
		if err != nil {
			return nil, err
		}
		if current != nil {
			locked = assignmentLocked(a, current.StartedAt)
			if locked {
				selectedTournamentId = a.TournamentId
				tableNumber = a.TableNumber
			}
		}
	}

	rows := make([]tournamentOption, len(list))
	for i, tournament := range list {
		rows[i] = tournamentOption{
			Id:       tournament.Id,
			Name:     tournament.Name,
			Status:   tournamentStatus(tournament),
			Selected: tournament.Id == selectedTournamentId,
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Name < rows[j].Name
	})

	tables, err := loadTables(s, selectedTournamentId)
	if err != nil {
		return nil, err
	}

	var assignedName string
	if selectedTournamentId != "" {
		for _, row := range rows {
			if row.Id == selectedTournamentId {
				assignedName = row.Name
				break
			}
		}
	}

	return fiber.Map{
		"Title":          "Konfiguracja",
		"Tournaments":    rows,
		"Tables":         tables,
		"TournamentId":   selectedTournamentId,
		"TournamentName": assignedName,
		"TableNumber":    tableNumber,
		"Locked":         locked,
		"Assigned":       hasAssignment,
	}, nil
}

func loadTables(s stores, tournamentId string) ([]tableOption, error) {
	if tournamentId == "" {
		return nil, nil
	}
	tables, err := (&tournament_query.ListTablesQuery{TournamentId: tournamentId}).Execute(s.tables)
	if err != nil {
		return nil, err
	}
	rows := make([]tableOption, len(tables))
	for i, table := range tables {
		rows[i] = tableOption{Number: table.Number}
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Number < rows[j].Number
	})
	return rows, nil
}

func waitingData(tournament *tournament_domain.TournamentDto, a assignment) fiber.Map {
	return fiber.Map{
		"Title":          "Oczekiwanie",
		"TournamentId":   tournament.Id,
		"TournamentName": tournament.Name,
		"TableNumber":    a.TableNumber,
	}
}

func findTournament(s stores, id string) (*tournament_domain.TournamentDto, error) {
	tournament, err := s.tournamentRead.FindById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if tournament == nil || tournament.Id == "" {
		return nil, nil
	}
	return tournament, nil
}

func tournamentStatus(t tournament_domain.TournamentDto) string {
	if t.FinishedAt != "" {
		return "zakończony"
	}
	if t.StartedAt != "" {
		return "w trakcie"
	}
	return "nie rozpoczęty"
}

type tournamentOption struct {
	Id       string
	Name     string
	Status   string
	Selected bool
}

type tableOption struct {
	Number int
}
