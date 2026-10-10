package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"bridge-tab/internal/idutil"
	rounds_domain "bridge-tab/internal/rounds-registration/domain"
	rounds_infra "bridge-tab/internal/rounds-registration/infrastructure"
	score_app "bridge-tab/internal/score-registration/application"
	score_domain "bridge-tab/internal/score-registration/domain"
	score_infra "bridge-tab/internal/score-registration/infrastructure"
	tournament_cmd "bridge-tab/internal/tournament-management/application/command"
	tournament_domain "bridge-tab/internal/tournament-management/domain"
	tournament_infra "bridge-tab/internal/tournament-management/infrastructure"
	user_domain "bridge-tab/internal/user/domain"
	user_infra "bridge-tab/internal/user/infrastructure"

	"github.com/gofiber/fiber/v2"
)

type stores struct {
	tournament     tournament_domain.TournamentRepository
	tournamentRead tournament_domain.TournamentReadRepository
	teams          tournament_domain.TeamReadRepository
	tables         tournament_domain.TableReadRepository
	boards         tournament_domain.BoardProtocolReadRepository
	sets           tournament_domain.SetReadRepository
	users          user_domain.UserReadRepository
	sessions       rounds_domain.GameSessionRepository
	sessionRead    rounds_domain.GameSessionReadRepository
	played         score_domain.PlayedResultRepository
	scores         score_domain.TournamentScoreRepository
	scoreRead      score_domain.TournamentScoreReadRepository
}

func newStores(c *fiber.Ctx) stores {
	tx := c.Locals("tx").(*sql.Tx)
	ctx := c.UserContext()
	return stores{
		tournament:     &tournament_infra.PostgresTournamentRepository{Ctx: ctx, Tx: tx},
		tournamentRead: &tournament_infra.PostgresTournamentReadRepository{Ctx: ctx, Tx: tx},
		teams:          &tournament_infra.PostgresTeamReadRepository{Ctx: ctx, Tx: tx},
		tables:         &tournament_infra.PostgresTableReadRepository{Ctx: ctx, Tx: tx},
		boards:         &tournament_infra.PostgresBoardProtocolReadRepository{Ctx: ctx, Tx: tx},
		sets:           &tournament_infra.PostgresSetReadRepository{Ctx: ctx, Tx: tx},
		users:          &user_infra.PostgresUserRepository{Ctx: ctx, Tx: tx},
		sessions:       &rounds_infra.PostgresGameSessionRepository{Ctx: ctx, Tx: tx},
		sessionRead:    &rounds_infra.PostgresGameSessionReadRepository{Ctx: ctx, Tx: tx},
		played:         &score_infra.PostgresPlayedResultRepository{Ctx: ctx, Tx: tx},
		scores:         &score_infra.PostgresTournamentScoreRepository{Ctx: ctx, Tx: tx},
		scoreRead:      &score_infra.PostgresTournamentScoreReadRepository{Ctx: ctx, Tx: tx},
	}
}

func view(c *fiber.Ctx, name string, data fiber.Map) error {
	return c.Render(name, prepare(c, data), "layout")
}

func reject(c *fiber.Ctx, name string, data fiber.Map, err error) error {
	if tx, ok := c.Locals("tx").(*sql.Tx); ok && c.Locals("txClosed") == nil {
		_ = tx.Rollback()
		c.Locals("txClosed", true)
	}
	if data == nil {
		data = fiber.Map{}
	}
	data["Error"] = explain(err).Error()
	return c.Status(fiber.StatusBadRequest).Render(name, prepare(c, data), "layout")
}

func prepare(c *fiber.Ctx, data fiber.Map) fiber.Map {
	if data == nil {
		data = fiber.Map{}
	}
	if _, ok := data["CSRF"]; !ok {
		token, err := csrfToken(c)
		if err != nil {
			token = ""
		}
		data["CSRF"] = token
	}
	defaults := fiber.Map{
		"Error":          "",
		"Title":          "",
		"TournamentId":   "",
		"TournamentName": "",
	}
	for key, value := range defaults {
		if _, ok := data[key]; !ok {
			data[key] = value
		}
	}
	return data
}

func pageData(t *tournament_domain.TournamentDto) fiber.Map {
	return fiber.Map{
		"Title":          t.Name,
		"TournamentId":   t.Id,
		"TournamentName": t.Name,
		"Started":        t.StartedAt != "",
		"Finished":       t.FinishedAt != "",
		"Status":         statusLabel(*t),
		"Code":           t.Id,
	}
}

func loadTournament(c *fiber.Ctx, s stores) (*tournament_domain.TournamentDto, error) {
	t, err := s.tournamentRead.FindById(c.Params("id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fiber.NewError(fiber.StatusNotFound, "turniej nie istnieje")
		}
		return nil, err
	}
	if t == nil || t.Id == "" {
		return nil, fiber.NewError(fiber.StatusNotFound, "turniej nie istnieje")
	}
	return t, nil
}

func statusLabel(t tournament_domain.TournamentDto) string {
	if t.FinishedAt != "" {
		return "zakończony"
	}
	if t.StartedAt != "" {
		return "w trakcie"
	}
	return "nie rozpoczęty"
}

func explain(err error) error {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return errors.New("nie znaleziono")
	case errors.Is(err, tournament_cmd.ErrRoundsNotAllPlayed):
		return errors.New("nie wszystkie rozdania zostały rozegrane")
	case errors.Is(err, score_app.ErrPlayedResultsNotFound):
		return errors.New("brak wyników do obliczenia — najpierw zakończ turniej")
	case errors.Is(err, tournament_domain.ErrTournamentAlreadyStarted):
		return errors.New("turniej jest już rozpoczęty")
	case errors.Is(err, tournament_domain.ErrSomeTeamHasNoMembers):
		return errors.New("każda drużyna musi mieć co najmniej jednego zawodnika")
	case errors.Is(err, tournament_domain.ErrTournamentNotStarted):
		return errors.New("turniej nie jest rozpoczęty")
	case errors.Is(err, tournament_domain.ErrTournamentAlreadyFinished):
		return errors.New("turniej jest już zakończony")
	case errors.Is(err, tournament_domain.ErrInvalidTeamName):
		return errors.New("nazwa drużyny nie może być pusta")
	case errors.Is(err, tournament_domain.ErrTeamNameAlreadyExists):
		return errors.New("drużyna o tej nazwie już istnieje")
	case errors.Is(err, tournament_domain.ErrTeamNumberAlreadyExists):
		return errors.New("numer drużyny jest zajęty")
	case errors.Is(err, tournament_domain.ErrTableNumberAlreadyExists):
		return errors.New("numer stołu jest zajęty")
	case errors.Is(err, tournament_domain.ErrTableReferencedByBoardProtocol):
		return errors.New("stół jest użyty w protokole")
	case errors.Is(err, tournament_domain.ErrBoardProtocolBelongsToSet):
		return errors.New("rozdanie należy do seta")
	case errors.Is(err, rounds_domain.ErrRoundAlreadyPlayed):
		return errors.New("to rozdanie jest już zapisane")
	case errors.Is(err, rounds_domain.ErrRoundNotPlayed):
		return errors.New("to rozdanie nie ma jeszcze zapisu")
	case errors.Is(err, rounds_domain.ErrRoundNotFound):
		return errors.New("nie ma takiego rozdania")
	default:
		return err
	}
}

func nameMap(s stores) (map[string]string, error) {
	users, err := s.users.FindAll()
	if err != nil {
		return nil, err
	}
	names := map[string]string{}
	for _, user := range users {
		names[strings.ToLower(user.Id)] = user.Name
	}
	return names, nil
}

func playerNames(members []tournament_domain.ContestantDto, names map[string]string) string {
	parts := make([]string, 0, len(members))
	for _, member := range members {
		id := string(member.Id)
		if name, ok := names[strings.ToLower(id)]; ok && name != "" {
			parts = append(parts, name)
			continue
		}
		parts = append(parts, id)
	}
	return strings.Join(parts, ", ")
}

func teamCaption(teams []tournament_domain.TeamDto, id string, number int) string {
	for _, team := range teams {
		if idutil.SameId(team.Id, id) {
			return fmt.Sprintf("%d %s", team.Number, team.Name)
		}
	}
	if number != 0 {
		return strconv.Itoa(number)
	}
	return id
}

func teamByNumber(teams []tournament_domain.TeamDto, number int) string {
	for _, team := range teams {
		if team.Number == number {
			return fmt.Sprintf("%d %s", team.Number, team.Name)
		}
	}
	return strconv.Itoa(number)
}

func vulLabel(value string) string {
	switch value {
	case "NS":
		return "NS"
	case "EW":
		return "EW"
	case "Both":
		return "obie"
	default:
		return "brak"
	}
}

func joinInts(nums []int) string {
	parts := make([]string, len(nums))
	for i, n := range nums {
		parts[i] = strconv.Itoa(n)
	}
	return strings.Join(parts, ", ")
}

func optionalPositive(raw string) (int, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, nil
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 {
		return 0, errors.New("numer musi być dodatni")
	}
	return n, nil
}

func pathInt(c *fiber.Ctx, name string) (int, error) {
	n, err := strconv.Atoi(c.Params(name))
	if err != nil || n < 1 {
		return 0, fiber.NewError(fiber.StatusBadRequest, "zły numer")
	}
	return n, nil
}

func formatScore(v float64) string {
	return strconv.FormatFloat(v, 'f', 1, 64)
}

func formatPct(v float64) string {
	return strconv.FormatFloat(v, 'f', 2, 64)
}
