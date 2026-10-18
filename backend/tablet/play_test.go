package main

import (
	"testing"

	play_query "bridge-tab/internal/play/application/query"
	tournament_domain "bridge-tab/internal/tournament-management/domain"
)

func TestBuildTablePlayPicksFirstIncompleteRound(t *testing.T) {
	play := buildTablePlay(3, false, []play_query.ScheduledMatchDto{
		{DealNo: 1, RoundNo: 1, TableNumber: 3, NsTeamNumber: 1, EwTeamNumber: 2, Contract: "1N"},
		{DealNo: 2, RoundNo: 1, TableNumber: 3, NsTeamNumber: 1, EwTeamNumber: 2, Contract: "Pass"},
		{DealNo: 3, RoundNo: 2, TableNumber: 3, NsTeamNumber: 4, EwTeamNumber: 5},
		{DealNo: 4, RoundNo: 2, TableNumber: 3, NsTeamNumber: 4, EwTeamNumber: 5},
		{DealNo: 9, RoundNo: 1, TableNumber: 1, NsTeamNumber: 8, EwTeamNumber: 9},
	}, 0)
	if play.RoundNo != 2 {
		t.Fatalf("round %d", play.RoundNo)
	}
	if play.NsTeamNumber != 4 || play.EwTeamNumber != 5 {
		t.Fatalf("teams %d %d", play.NsTeamNumber, play.EwTeamNumber)
	}
	if play.AllDone || play.Empty || play.RoundComplete {
		t.Fatalf("done=%v empty=%v complete=%v", play.AllDone, play.Empty, play.RoundComplete)
	}
	if len(play.Boards) != 2 {
		t.Fatalf("boards %d", len(play.Boards))
	}
}

func TestBuildTablePlayHoldsCompletedRound(t *testing.T) {
	matches := []play_query.ScheduledMatchDto{
		{DealNo: 1, RoundNo: 1, TableNumber: 3, NsTeamNumber: 1, EwTeamNumber: 2, Contract: "1N"},
		{DealNo: 2, RoundNo: 1, TableNumber: 3, NsTeamNumber: 1, EwTeamNumber: 2, Contract: "Pass"},
		{DealNo: 3, RoundNo: 2, TableNumber: 3, NsTeamNumber: 4, EwTeamNumber: 5},
	}
	play := buildTablePlay(3, false, matches, 1)
	if play.RoundNo != 1 {
		t.Fatalf("round %d", play.RoundNo)
	}
	if !play.RoundComplete || !play.HasNextRound || play.NextRoundNo != 2 {
		t.Fatalf("complete=%v next=%v no=%d", play.RoundComplete, play.HasNextRound, play.NextRoundNo)
	}
	if play.NsTeamNumber != 1 {
		t.Fatalf("ns %d", play.NsTeamNumber)
	}
}

func TestBuildTablePlayAllDoneUsesLastRound(t *testing.T) {
	play := buildTablePlay(3, false, []play_query.ScheduledMatchDto{
		{DealNo: 1, RoundNo: 1, TableNumber: 3, NsTeamNumber: 1, EwTeamNumber: 2, Contract: "1N"},
		{DealNo: 3, RoundNo: 2, TableNumber: 3, NsTeamNumber: 4, EwTeamNumber: 5, Contract: "Pass"},
	}, 0)
	if !play.AllDone {
		t.Fatal("expected all done")
	}
	if play.RoundNo != 2 {
		t.Fatalf("round %d", play.RoundNo)
	}
	if play.HasNextRound {
		t.Fatal("no next round")
	}
}

func TestDisplayContract(t *testing.T) {
	cases := map[string]string{
		"":     "",
		"Pass": "Pass",
		"4H":   "4H",
		"3N":   "3NT",
		"3Nx":  "3NTx",
		"3Nxx": "3NTxx",
		"3NT":  "3NT",
		"3NTx": "3NTx",
	}
	for in, want := range cases {
		if got := displayContract(in); got != want {
			t.Errorf("displayContract(%q)=%q, want %q", in, got, want)
		}
	}
}

func TestFormatCallHTML(t *testing.T) {
	got := string(formatCallHTML("3NT"))
	if got != "3NT" {
		t.Fatalf("3NT -> %q", got)
	}
	got = string(formatCallHTML("4H"))
	if got != `4<span class="suit red">♥</span>` {
		t.Fatalf("4H -> %q", got)
	}
	got = string(formatLeadHTML("AH"))
	if got != `A<span class="suit red">♥</span>` {
		t.Fatalf("AH -> %q", got)
	}
}

func TestBuildTablePlayEmpty(t *testing.T) {
	play := buildTablePlay(9, false, []play_query.ScheduledMatchDto{
		{DealNo: 1, RoundNo: 1, TableNumber: 3, NsTeamNumber: 1, EwTeamNumber: 2},
	}, 0)
	if !play.Empty {
		t.Fatal("expected empty")
	}
}

func TestPlayForDeal(t *testing.T) {
	play, ok := playForDeal(3, false, []play_query.ScheduledMatchDto{
		{DealNo: 1, RoundNo: 1, TableNumber: 3, NsTeamNumber: 1, EwTeamNumber: 2, Contract: "1N"},
		{DealNo: 3, RoundNo: 2, TableNumber: 3, NsTeamNumber: 4, EwTeamNumber: 5},
	}, 1)
	if !ok || play.RoundNo != 1 || play.NsTeamNumber != 1 {
		t.Fatalf("got round %d ns %d ok %v", play.RoundNo, play.NsTeamNumber, ok)
	}
}

func TestVulLabel(t *testing.T) {
	cases := map[string]string{
		"":     "brak",
		"None": "brak",
		"NS":   "NS",
		"EW":   "EW",
		"Both": "obie",
	}
	for in, want := range cases {
		if got := vulLabel(in); got != want {
			t.Errorf("vulLabel(%q)=%q, want %q", in, got, want)
		}
	}
}

func TestToBoardViewsVulnerability(t *testing.T) {
	boards := toBoardViews([]play_query.ScheduledMatchDto{
		{DealNo: 1, Vulnerable: "NS"},
		{DealNo: 2, Vulnerable: "EW"},
		{DealNo: 3, Vulnerable: "Both"},
		{DealNo: 4, Vulnerable: "None"},
	})
	if !boards[0].NsVulnerable || boards[0].EwVulnerable || boards[0].VulnerableLabel != "NS" {
		t.Fatalf("NS board %+v", boards[0])
	}
	if boards[1].NsVulnerable || !boards[1].EwVulnerable || boards[1].VulnerableLabel != "EW" {
		t.Fatalf("EW board %+v", boards[1])
	}
	if !boards[2].NsVulnerable || !boards[2].EwVulnerable || boards[2].VulnerableLabel != "obie" {
		t.Fatalf("Both board %+v", boards[2])
	}
	if boards[3].NsVulnerable || boards[3].EwVulnerable || boards[3].VulnerableLabel != "brak" {
		t.Fatalf("None board %+v", boards[3])
	}
}

func TestTeamLabel(t *testing.T) {
	teams := []tournament_domain.TeamDto{
		{Number: 1, Name: "Mostostal"},
		{Number: 2, Name: "  "},
		{Number: 3, Name: ""},
	}
	if got := teamLabel(teams, 1); got != "Mostostal" {
		t.Fatalf("named %q", got)
	}
	if got := teamLabel(teams, 2); got != "Drużyna 2" {
		t.Fatalf("blank %q", got)
	}
	if got := teamLabel(teams, 3); got != "Drużyna 3" {
		t.Fatalf("empty %q", got)
	}
	if got := teamLabel(teams, 9); got != "Drużyna 9" {
		t.Fatalf("missing %q", got)
	}
	if got := teamLabel(teams, 0); got != "" {
		t.Fatalf("zero %q", got)
	}
}
