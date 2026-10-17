package main

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	play_app "bridge-tab/internal/play/application"
	play_query "bridge-tab/internal/play/application/query"
	play_domain "bridge-tab/internal/play/domain"
	tournament_query "bridge-tab/internal/tournament-management/application/query"
	tournament_domain "bridge-tab/internal/tournament-management/domain"

	"github.com/gofiber/fiber/v2"
)

type roundRow struct {
	RoundNo     int
	TableNumber int
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
	RoundNo     int
	TableNumber int
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

	playSessionId := tournament.Id
	if form.Played {
		cmd := play_app.EditBoardResultCommand{
			PlaySessionId:    playSessionId,
			FirstTeamNumber:  form.NS,
			SecondTeamNumber: form.EW,
			DealNo:           form.DealNo,
			Contract:         contract,
			Tricks:           tricks,
			Declarer:         declarer,
			OpeningLead:      lead,
		}
		if err := cmd.Execute(s.boardResult, s.playSession, s.played); err != nil {
			return reject(c, "round", data, err)
		}
	} else {
		cmd := play_app.RecordBoardResultByTeamsCommand{
			PlaySessionId:    playSessionId,
			FirstTeamNumber:  form.NS,
			SecondTeamNumber: form.EW,
			DealNo:           form.DealNo,
			Contract:         contract,
			Tricks:           tricks,
			Declarer:         declarer,
			OpeningLead:      lead,
		}
		if err := cmd.Execute(s.boardResult, s.playSession, s.played); err != nil {
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
	playSessionId := tournament.Id
	matches, err := (&play_query.ListMatchesQuery{PlaySessionId: playSessionId}).Execute(s.playSession, s.boardResults)
	if err != nil {
		return nil, err
	}
	showScores := tournament.FinishedAt != "" || reveal
	rows := make([]roundRow, len(matches))
	for i, match := range matches {
		tricks := ""
		if match.Contract != "" && match.Contract != "Pass" {
			tricks = strconv.Itoa(match.Tricks)
		}
		rows[i] = roundRow{
			RoundNo:     match.RoundNo,
			TableNumber: match.TableNumber,
			DealNo:      match.DealNo,
			NS:          match.NsTeamNumber,
			EW:          match.EwTeamNumber,
			Played:      match.Contract != "",
			Contract:    match.Contract,
			Declarer:    match.Declarer,
			Tricks:      tricks,
			OpeningLead: match.OpeningLead,
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].RoundNo != rows[j].RoundNo {
			return rows[i].RoundNo < rows[j].RoundNo
		}
		if rows[i].TableNumber != rows[j].TableNumber {
			return rows[i].TableNumber < rows[j].TableNumber
		}
		if rows[i].DealNo != rows[j].DealNo {
			return rows[i].DealNo < rows[j].DealNo
		}
		if rows[i].NS != rows[j].NS {
			return rows[i].NS < rows[j].NS
		}
		return rows[i].EW < rows[j].EW
	})
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
	playSessionId := tournament.Id
	match, err := (&play_query.GetMatchByTeamsQuery{
		PlaySessionId:    playSessionId,
		FirstTeamNumber:  ns,
		SecondTeamNumber: ew,
		DealNo:           deal,
	}).Execute(s.playSession, s.boardResults)
	if err != nil {
		if errors.Is(err, play_domain.ErrMatchNotFound) || errors.Is(err, play_domain.ErrTeamNotFound) || errors.Is(err, play_domain.ErrPlaySessionNotProvisioned) {
			return form, nil
		}
		return roundForm{}, err
	}
	form.Found = true
	form.Played = match.Contract != ""
	form.Contract = match.Contract
	form.Declarer = match.Declarer
	form.OpeningLead = match.OpeningLead
	form.RoundNo = match.RoundNo
	form.TableNumber = match.TableNumber
	form.NS = match.NsTeamNumber
	form.EW = match.EwTeamNumber
	form.NSLabel = teamByNumber(teams, match.NsTeamNumber)
	form.EWLabel = teamByNumber(teams, match.EwTeamNumber)
	if match.Contract != "" && match.Contract != "Pass" {
		form.Tricks = strconv.Itoa(match.Tricks)
	}
	return form, nil
}

func postedScore(c *fiber.Ctx) (string, string, string, int, error) {
	contract := normalizePostedContract(c.FormValue("contract"))
	declarer := strings.ToUpper(strings.TrimSpace(c.FormValue("declarer")))
	lead := strings.ToUpper(strings.TrimSpace(c.FormValue("openingLead")))
	rawTricks := strings.TrimSpace(c.FormValue("tricks"))
	if rawTricks == "" {
		if contract != "Pass" {
			return "", "", "", 0, errors.New("tricks is empty")
		}
		return contract, declarer, lead, 0, nil
	}
	tricks, err := strconv.Atoi(rawTricks)
	if err != nil {
		return "", "", "", 0, errors.New("zła liczba lew")
	}
	return contract, declarer, lead, tricks, nil
}
