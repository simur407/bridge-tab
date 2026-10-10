package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	rounds_cmd "bridge-tab/internal/rounds-registration/application"
	rounds_query "bridge-tab/internal/rounds-registration/application/query"
	tournament_query "bridge-tab/internal/tournament-management/application/query"
	tournament_domain "bridge-tab/internal/tournament-management/domain"

	"github.com/gofiber/fiber/v2"
)

type roundRow struct {
	DealNo      int
	NS          int
	EW          int
	Played      bool
	Contract    string
	Declarer    string
	Tricks      string
	OpeningLead string
}

type roundForm struct {
	DealNo      int
	NS          int
	EW          int
	NSLabel     string
	EWLabel     string
	Contract    string
	Declarer    string
	Tricks      string
	OpeningLead string
	Played      bool
	Found       bool
}

func getRounds(c *fiber.Ctx) error {
	s := newStores(c)
	tournament, err := loadTournament(c, s)
	if err != nil {
		return err
	}
	data, err := roundsPage(s, tournament, c.Query("reveal") == "1")
	if err != nil {
		return err
	}
	return view(c, "rounds", data)
}

func getRound(c *fiber.Ctx) error {
	s := newStores(c)
	tournament, err := loadTournament(c, s)
	if err != nil {
		return err
	}
	form, err := loadRoundForm(c, s, tournament)
	if err != nil {
		return err
	}
	data := pageData(tournament)
	data["Round"] = form
	data["Reveal"] = c.Query("reveal") == "1"
	if !form.Found {
		data["Error"] = "rozdanie nie istnieje"
	}
	return view(c, "round", data)
}

func postRound(c *fiber.Ctx) error {
	s := newStores(c)
	tournament, err := loadTournament(c, s)
	if err != nil {
		return err
	}
	form, err := loadRoundForm(c, s, tournament)
	if err != nil {
		return err
	}
	data := pageData(tournament)
	data["Reveal"] = c.FormValue("reveal") == "1"
	if !form.Found {
		data["Round"] = form
		return reject(c, "round", data, errors.New("rozdanie nie istnieje"))
	}

	contract, declarer, lead, tricks, err := postedScore(c)
	if err != nil {
		form.Contract = c.FormValue("contract")
		form.Declarer = c.FormValue("declarer")
		form.OpeningLead = c.FormValue("openingLead")
		form.Tricks = c.FormValue("tricks")
		data["Round"] = form
		return reject(c, "round", data, err)
	}
	form.Contract = contract
	form.Declarer = declarer
	form.OpeningLead = lead
	form.Tricks = strings.TrimSpace(c.FormValue("tricks"))
	data["Round"] = form

	if form.Played {
		cmd := rounds_cmd.EditRoundCommand{
			GameSessionId:    tournament.Id,
			FirstTeamNumber:  form.NS,
			SecondTeamNumber: form.EW,
			DealNo:           form.DealNo,
			Contract:         contract,
			Tricks:           tricks,
			Declarer:         declarer,
			OpeningLead:      lead,
		}
		if err := cmd.Execute(s.sessions, s.teams); err != nil {
			return reject(c, "round", data, err)
		}
	} else {
		cmd := rounds_cmd.RecordRoundCommand{
			GameSessionId:    tournament.Id,
			FirstTeamNumber:  form.NS,
			SecondTeamNumber: form.EW,
			DealNo:           form.DealNo,
			Contract:         contract,
			Tricks:           tricks,
			Declarer:         declarer,
			OpeningLead:      lead,
		}
		if err := cmd.Execute(s.sessions, s.teams); err != nil {
			return reject(c, "round", data, err)
		}
	}

	target := fmt.Sprintf("/tournaments/%s/rounds/%d/%d/%d", tournament.Id, form.DealNo, form.NS, form.EW)
	if c.FormValue("reveal") == "1" {
		target += "?reveal=1"
	}
	return c.Redirect(target, fiber.StatusSeeOther)
}

func roundsPage(s stores, tournament *tournament_domain.TournamentDto, reveal bool) (fiber.Map, error) {
	rounds, err := (&rounds_query.ListRoundsQuery{GameSessionId: tournament.Id}).Execute(s.sessionRead)
	if err != nil {
		return nil, err
	}
	showScores := tournament.FinishedAt != "" || reveal
	rows := make([]roundRow, len(rounds))
	for i, round := range rounds {
		tricks := ""
		if round.Contract != "" && round.Contract != "Pass" {
			tricks = strconv.Itoa(round.Tricks)
		}
		rows[i] = roundRow{
			DealNo:      round.DealNo,
			NS:          round.NsTeamNumber,
			EW:          round.EwTeamNumber,
			Played:      round.Contract != "",
			Contract:    round.Contract,
			Declarer:    round.Declarer,
			Tricks:      tricks,
			OpeningLead: round.OpeningLead,
		}
	}
	data := pageData(tournament)
	data["Rounds"] = rows
	data["ShowScores"] = showScores
	data["CanHide"] = tournament.FinishedAt == ""
	return data, nil
}

func loadRoundForm(c *fiber.Ctx, s stores, tournament *tournament_domain.TournamentDto) (roundForm, error) {
	deal, err := pathInt(c, "deal")
	if err != nil {
		return roundForm{}, err
	}
	ns, err := pathInt(c, "ns")
	if err != nil {
		return roundForm{}, err
	}
	ew, err := pathInt(c, "ew")
	if err != nil {
		return roundForm{}, err
	}
	teams, err := (&tournament_query.ListTeamsQuery{TournamentId: tournament.Id}).Execute(s.teams)
	if err != nil {
		return roundForm{}, err
	}
	form := roundForm{
		DealNo:  deal,
		NS:      ns,
		EW:      ew,
		NSLabel: teamByNumber(teams, ns),
		EWLabel: teamByNumber(teams, ew),
	}
	round, err := (&rounds_query.GetRoundByTeamsQuery{
		GameSessionId:    tournament.Id,
		FirstTeamNumber:  ns,
		SecondTeamNumber: ew,
		DealNo:           deal,
	}).Execute(s.sessionRead, s.teams)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return form, nil
		}
		return roundForm{}, err
	}
	if round == nil {
		return form, nil
	}
	form.Found = true
	form.Played = round.Contract != ""
	form.Contract = round.Contract
	form.Declarer = round.Declarer
	form.OpeningLead = round.OpeningLead
	form.NS = round.NsTeamNumber
	form.EW = round.EwTeamNumber
	form.NSLabel = teamByNumber(teams, round.NsTeamNumber)
	form.EWLabel = teamByNumber(teams, round.EwTeamNumber)
	if round.Contract != "" && round.Contract != "Pass" {
		form.Tricks = strconv.Itoa(round.Tricks)
	}
	return form, nil
}

func postedScore(c *fiber.Ctx) (string, string, string, int, error) {
	contract := normalizePostedContract(c.FormValue("contract"))
	declarer := strings.ToUpper(strings.TrimSpace(c.FormValue("declarer")))
	lead := strings.ToUpper(strings.TrimSpace(c.FormValue("openingLead")))
	rawTricks := strings.TrimSpace(c.FormValue("tricks"))
	if rawTricks == "" {
		return contract, declarer, lead, 0, nil
	}
	tricks, err := strconv.Atoi(rawTricks)
	if err != nil {
		return "", "", "", 0, errors.New("zła liczba lew")
	}
	return contract, declarer, lead, tricks, nil
}
