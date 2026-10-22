package main

import (
	"database/sql"
	"net/url"
	"strings"
	"time"

	tournament_domain "bridge-tab/internal/tournament-management/domain"
	umpirecall "bridge-tab/internal/umpire-call"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var umpireZone = func() *time.Location {
	loc, err := time.LoadLocation("Europe/Warsaw")
	if err != nil {
		return time.Local
	}
	return loc
}()

type umpireRow struct {
	Id          string
	TableNumber int
	RoundNo     int
	Ns          string
	Ew          string
	Source      string
	When        string
	Answered    string
}

func getUmpireCalls(c *fiber.Ctx) error {
	s := newStores(c)
	tournament, err := loadTournament(c, s)
	if err != nil {
		return err
	}
	calls, err := umpireRepo(c).List(tournament.Id)
	if err != nil {
		return err
	}
	teams, err := s.teams.FindAll(&tournament.Id)
	if err != nil {
		return err
	}
	open := make([]umpireRow, 0)
	answered := make([]umpireRow, 0)
	for _, call := range calls {
		row := umpireRowFrom(call, teams)
		if call.AnsweredAt == nil {
			open = append(open, row)
			continue
		}
		answered = append(answered, row)
	}
	data := pageData(tournament)
	data["Open"] = open
	data["Answered"] = answered
	return view(c, "umpire-calls", data)
}

func getOpenUmpireCalls(c *fiber.Ctx) error {
	s := newStores(c)
	tournament, err := loadTournament(c, s)
	if err != nil {
		return err
	}
	calls, err := umpireRepo(c).ListOpen(tournament.Id)
	if err != nil {
		return err
	}
	teams, err := s.teams.FindAll(&tournament.Id)
	if err != nil {
		return err
	}
	rows := make([]umpireRow, 0, len(calls))
	for _, call := range calls {
		rows = append(rows, umpireRowFrom(call, teams))
	}
	data := fiber.Map{
		"TournamentId": tournament.Id,
		"Calls":        rows,
		"OpenCount":    len(rows),
	}
	return c.Render("umpire-open", prepare(c, data))
}

func postAnswerUmpireCall(c *fiber.Ctx) error {
	tournamentId := c.Params("id")
	callId := c.Params("callId")
	if err := uuid.Validate(callId); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "zły identyfikator")
	}
	s := newStores(c)
	if _, err := loadTournament(c, s); err != nil {
		return err
	}
	if _, err := umpireRepo(c).MarkAnswered(tournamentId, callId); err != nil {
		return err
	}
	return c.Redirect(adminReturn(c, "/tournaments/"+tournamentId+"/umpire-calls"), fiber.StatusSeeOther)
}

func umpireRepo(c *fiber.Ctx) *umpirecall.Repository {
	return &umpirecall.Repository{
		Ctx: c.UserContext(),
		Tx:  c.Locals("tx").(*sql.Tx),
	}
}

func umpireRowFrom(call umpirecall.Call, teams []tournament_domain.TeamDto) umpireRow {
	row := umpireRow{
		Id:          call.Id,
		TableNumber: call.TableNumber,
		RoundNo:     call.RoundNo,
		Ns:          strings.TrimSpace(teamByNumber(teams, call.NsTeamNumber)),
		Ew:          strings.TrimSpace(teamByNumber(teams, call.EwTeamNumber)),
		Source:      umpireSource(call.Source),
		When:        call.CreatedAt.In(umpireZone).Format("15:04:05"),
	}
	if call.AnsweredAt != nil {
		row.Answered = call.AnsweredAt.In(umpireZone).Format("15:04:05")
	}
	return row
}

func umpireSource(source string) string {
	if source == umpirecall.SourceContestant {
		return "zawodnik"
	}
	return "tablet"
}

func adminReturn(c *fiber.Ctx, fallback string) string {
	ref := c.Get("Referer")
	u, err := url.Parse(ref)
	if err != nil || u.Path == "" || !strings.HasPrefix(u.Path, "/") || strings.HasPrefix(u.Path, "//") {
		return fallback
	}
	if u.Host != "" && !strings.EqualFold(u.Hostname(), c.Hostname()) {
		return fallback
	}
	return u.RequestURI()
}
