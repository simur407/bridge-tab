package main

import (
	"database/sql"
	"errors"

	play_domain "bridge-tab/internal/play/domain"
	play_infra "bridge-tab/internal/play/infrastructure"
	score_domain "bridge-tab/internal/score-registration/domain"
	score_infra "bridge-tab/internal/score-registration/infrastructure"
	tournament_domain "bridge-tab/internal/tournament-management/domain"
	tournament_infra "bridge-tab/internal/tournament-management/infrastructure"

	"github.com/gofiber/fiber/v2"
)

type stores struct {
	tournamentRead tournament_domain.TournamentReadRepository
	tables         tournament_domain.TableReadRepository
	teams          tournament_domain.TeamReadRepository
	playSession    play_domain.PlaySessionRepository
	boardResult    play_domain.BoardResultRepository
	boardResults   play_domain.BoardResultReadRepository
	played         score_domain.PlayedResultRepository
}

func newStores(c *fiber.Ctx) stores {
	tx := c.Locals("tx").(*sql.Tx)
	ctx := c.UserContext()
	return stores{
		tournamentRead: &tournament_infra.PostgresTournamentReadRepository{Ctx: ctx, Tx: tx},
		tables:         &tournament_infra.PostgresTableReadRepository{Ctx: ctx, Tx: tx},
		teams:          &tournament_infra.PostgresTeamReadRepository{Ctx: ctx, Tx: tx},
		playSession:    &play_infra.PostgresPlaySessionRepository{Ctx: ctx, Tx: tx},
		boardResult:    &play_infra.PostgresBoardResultRepository{Ctx: ctx, Tx: tx},
		boardResults:   &play_infra.PostgresBoardResultReadRepository{Ctx: ctx, Tx: tx},
		played:         &score_infra.PostgresPlayedResultRepository{Ctx: ctx, Tx: tx},
	}
}

func view(c *fiber.Ctx, name string, data fiber.Map) error {
	return c.Render(name, prepare(c, data), "layout")
}

func playView(c *fiber.Ctx, name string, data fiber.Map) error {
	return c.Render(name, prepare(c, data), "play-layout")
}

func fragment(c *fiber.Ctx, name string, data fiber.Map) error {
	return c.Render(name, prepare(c, data))
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

func rejectFragment(c *fiber.Ctx, name string, data fiber.Map, err error) error {
	if tx, ok := c.Locals("tx").(*sql.Tx); ok && c.Locals("txClosed") == nil {
		_ = tx.Rollback()
		c.Locals("txClosed", true)
	}
	if data == nil {
		data = fiber.Map{}
	}
	data["Error"] = explain(err).Error()
	return c.Status(fiber.StatusBadRequest).Render(name, prepare(c, data))
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
		"Error":   "",
		"Title":   "",
		"Success": "",
	}
	for key, value := range defaults {
		if _, ok := data[key]; !ok {
			data[key] = value
		}
	}
	return data
}

func explain(err error) error {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return errors.New("nie znaleziono")
	case errors.Is(err, play_domain.ErrPlaySessionNotProvisioned):
		return errors.New("rozgrywka nie została rozpoczęta")
	case errors.Is(err, play_domain.ErrMatchNotFound):
		return errors.New("nie ma takiego rozdania")
	case errors.Is(err, play_domain.ErrBoardResultAlreadyRecorded):
		return errors.New("rozdanie już zostało wpisane")
	case errors.Is(err, play_domain.ErrBoardResultNotFound):
		return errors.New("nie znaleziono rozdania")
	case errors.Is(err, play_domain.ErrTeamNotFound):
		return errors.New("nie znaleziono zespołu")
	case errors.Is(err, tournament_domain.ErrNoSuchTableInTournament):
		return errors.New("nie znaleziono stołu")
	default:
		if err == nil {
			return errors.New("błąd")
		}
		return err
	}
}
