package main

import (
	"bytes"
	"strings"
	"testing"

	play_query "bridge-tab/internal/play/application/query"
	tournament_domain "bridge-tab/internal/tournament-management/domain"

	"github.com/gofiber/template/html/v2"
)

func TestBuildTeamPlayPicksFirstIncompleteRound(t *testing.T) {
	play := buildTeamPlay(2, false, []play_query.ScheduledMatchDto{
		{DealNo: 1, RoundNo: 1, TableNumber: 3, NsTeamNumber: 1, EwTeamNumber: 2, Contract: "1N"},
		{DealNo: 2, RoundNo: 1, TableNumber: 3, NsTeamNumber: 1, EwTeamNumber: 2, Contract: "Pass"},
		{DealNo: 3, RoundNo: 2, TableNumber: 4, NsTeamNumber: 2, EwTeamNumber: 5},
		{DealNo: 4, RoundNo: 2, TableNumber: 4, NsTeamNumber: 2, EwTeamNumber: 5},
		{DealNo: 9, RoundNo: 1, TableNumber: 1, NsTeamNumber: 8, EwTeamNumber: 9},
	}, 0)
	if play.RoundNo != 2 || play.TableNumber != 4 {
		t.Fatalf("round %d table %d", play.RoundNo, play.TableNumber)
	}
	if play.Direction != "NS" || play.OpponentNumber != 5 {
		t.Fatalf("seat %s vs %d", play.Direction, play.OpponentNumber)
	}
	if play.AllDone || play.Empty || play.RoundComplete {
		t.Fatalf("done=%v empty=%v complete=%v", play.AllDone, play.Empty, play.RoundComplete)
	}
	if play.HasNextRound {
		t.Fatal("round 2 is the last round")
	}
	if len(play.Boards) != 2 || play.Boards[0].DealNo != 3 {
		t.Fatalf("boards %+v", play.Boards)
	}
}

func TestBuildTeamPlayShowsNextSeatingDuringRound(t *testing.T) {
	teams := []tournament_domain.TeamDto{
		{Number: 2, Name: "Alfa"},
		{Number: 5, Name: "Beta"},
		{Number: 8, Name: "Gamma"},
	}
	play := labelTeamPlay(buildTeamPlay(2, false, []play_query.ScheduledMatchDto{
		{DealNo: 1, RoundNo: 1, TableNumber: 3, NsTeamNumber: 5, EwTeamNumber: 2},
		{DealNo: 2, RoundNo: 1, TableNumber: 3, NsTeamNumber: 5, EwTeamNumber: 2},
		{DealNo: 7, RoundNo: 2, TableNumber: 6, NsTeamNumber: 2, EwTeamNumber: 8},
	}, 0), teams)
	if play.RoundNo != 1 || play.Direction != "EW" || play.OpponentLabel != "5 Beta" {
		t.Fatalf("current %+v", play)
	}
	if play.TeamLabel != "2 Alfa" {
		t.Fatalf("team %q", play.TeamLabel)
	}
	if !play.HasNextRound || play.NextRoundNo != 2 || play.NextTableNumber != 6 || play.NextDirection != "NS" || play.NextOpponentLabel != "8 Gamma" {
		t.Fatalf("next round=%d table=%d dir=%s opp=%q has=%v", play.NextRoundNo, play.NextTableNumber, play.NextDirection, play.NextOpponentLabel, play.HasNextRound)
	}
}

func TestBuildTeamPlayHoldsCompletedRound(t *testing.T) {
	play := buildTeamPlay(1, false, []play_query.ScheduledMatchDto{
		{DealNo: 1, RoundNo: 1, TableNumber: 3, NsTeamNumber: 1, EwTeamNumber: 2, Contract: "1N"},
		{DealNo: 2, RoundNo: 1, TableNumber: 3, NsTeamNumber: 1, EwTeamNumber: 2, Contract: "Pass"},
		{DealNo: 3, RoundNo: 2, TableNumber: 4, NsTeamNumber: 1, EwTeamNumber: 5},
	}, 1)
	if play.RoundNo != 1 || !play.Held || !play.RoundComplete {
		t.Fatalf("round=%d held=%v complete=%v", play.RoundNo, play.Held, play.RoundComplete)
	}
	if !play.HasNextRound || play.NextRoundNo != 2 || play.NextTableNumber != 4 || play.NextOpponentNumber != 5 {
		t.Fatalf("next %+v", play)
	}
}

func TestBuildTeamPlayAllDoneUsesLastRound(t *testing.T) {
	play := buildTeamPlay(1, false, []play_query.ScheduledMatchDto{
		{DealNo: 1, RoundNo: 1, TableNumber: 3, NsTeamNumber: 1, EwTeamNumber: 2, Contract: "1N"},
		{DealNo: 3, RoundNo: 2, TableNumber: 4, NsTeamNumber: 4, EwTeamNumber: 1, Contract: "Pass"},
	}, 0)
	if !play.AllDone || play.RoundNo != 2 || play.HasNextRound || play.Held {
		t.Fatalf("all=%v round=%d next=%v held=%v", play.AllDone, play.RoundNo, play.HasNextRound, play.Held)
	}
	if play.Direction != "EW" || play.TableNumber != 4 || play.OpponentNumber != 4 {
		t.Fatalf("seat %s table %d vs %d", play.Direction, play.TableNumber, play.OpponentNumber)
	}
}

func TestBuildTeamPlayEmpty(t *testing.T) {
	play := buildTeamPlay(9, false, []play_query.ScheduledMatchDto{
		{DealNo: 1, RoundNo: 1, TableNumber: 3, NsTeamNumber: 1, EwTeamNumber: 2},
	}, 0)
	if !play.Empty {
		t.Fatal("expected empty")
	}
}

func TestTeamPlayForDeal(t *testing.T) {
	play, ok := teamPlayForDeal(2, false, []play_query.ScheduledMatchDto{
		{DealNo: 1, RoundNo: 1, TableNumber: 3, NsTeamNumber: 1, EwTeamNumber: 2, Contract: "1N"},
		{DealNo: 3, RoundNo: 2, TableNumber: 6, NsTeamNumber: 2, EwTeamNumber: 8},
	}, 1)
	if !ok || play.RoundNo != 1 || play.Direction != "EW" || play.TableNumber != 3 {
		t.Fatalf("got round %d dir %s table %d ok %v", play.RoundNo, play.Direction, play.TableNumber, ok)
	}
	if _, ok := teamPlayForDeal(2, false, []play_query.ScheduledMatchDto{
		{DealNo: 1, RoundNo: 1, TableNumber: 3, NsTeamNumber: 1, EwTeamNumber: 4},
	}, 1); ok {
		t.Fatal("expected miss")
	}
}

func TestTeamPlayFormatters(t *testing.T) {
	if got := displayContract("3N"); got != "3NT" {
		t.Fatalf("display %q", got)
	}
	if got := string(formatCallHTML("4H")); got != `4<span class="suit red">♥</span>` {
		t.Fatalf("call %q", got)
	}
	if got := string(formatLeadHTML("10S")); got != `10<span class="suit club">♠</span>` {
		t.Fatalf("lead %q", got)
	}
	if got := string(formatCallHTML("Pass")); got != "Pas" {
		t.Fatalf("pass %q", got)
	}
	boards := toBoardViews([]play_query.ScheduledMatchDto{{DealNo: 1, Vulnerable: "Both", Contract: "2H"}})
	if !boards[0].Recorded || !boards[0].NsVulnerable || !boards[0].EwVulnerable || boards[0].VulnerableLabel != "obie" {
		t.Fatalf("board %+v", boards[0])
	}
}

func TestTeamLabel(t *testing.T) {
	teams := []tournament_domain.TeamDto{
		{Number: 1, Name: "Mostostal"},
		{Number: 2, Name: "  "},
	}
	if got := teamLabel(teams, 1); got != "1 Mostostal" {
		t.Fatalf("named %q", got)
	}
	if got := teamLabel(teams, 2); got != "Drużyna 2" {
		t.Fatalf("blank %q", got)
	}
	if got := teamLabel(teams, 0); got != "" {
		t.Fatalf("zero %q", got)
	}
}

func TestFrontendTemplatesParse(t *testing.T) {
	engine := html.New("frontend", ".html")
	engine.AddFunc("base", func() string { return "" })
	if err := engine.Load(); err != nil {
		t.Fatal(err)
	}
	data := fiberMapForPlay(false)
	var page bytes.Buffer
	if err := engine.Render(&page, "team-play", data, "layout"); err != nil {
		t.Fatal(err)
	}
	body := page.String()
	if !strings.Contains(body, "Runda 1") || !strings.Contains(body, "nierozgrane") {
		t.Fatalf("round missing from page: %s", body)
	}
	if strings.Contains(body, `hx-target="#team-play"`) || strings.Contains(body, `href="/round-registration/abc/boards/`) {
		t.Fatalf("round panel is still clickable or a swap target: %s", body)
	}
	var round bytes.Buffer
	if err := engine.Render(&round, "team-play-round", data); err != nil {
		t.Fatal(err)
	}
	fragment := round.String()
	if !strings.HasPrefix(strings.TrimSpace(fragment), `<div id="team-play"`) || strings.Contains(fragment, "Wezwij sędziego") {
		t.Fatalf("poll fragment: %s", fragment)
	}
	data["CanReadScores"] = true
	data["CanRegister"] = true
	var linked bytes.Buffer
	if err := engine.Render(&linked, "team-play-round", data); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(linked.String(), `href="/round-registration/abc/boards/1"`) {
		t.Fatalf("register mode should link boards: %s", linked.String())
	}
}

func fiberMapForPlay(read bool) map[string]any {
	return map[string]any{
		"Title":             "Turniej",
		"GameSessionId":     "abc",
		"TournamentName":    "Turniej",
		"CanCallUmpire":     true,
		"UserId":            "anna",
		"TeamLabel":         "1 Mostostal",
		"RoundNo":           1,
		"TableNumber":       3,
		"Direction":         "NS",
		"OpponentLabel":     "2 Wisła",
		"Boards":            []boardView{{DealNo: 1, VulnerableLabel: "brak"}},
		"HasNextRound":      true,
		"NextRoundNo":       2,
		"NextTableNumber":   4,
		"NextDirection":     "EW",
		"NextOpponentLabel": "5 Czarni",
		"CanRegister":       read,
		"CanReadScores":     read,
	}
}
