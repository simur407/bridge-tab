package score_registration_test

import (
	"testing"

	. "bridge-tab/internal/score-registration/domain"
)

func TestMatchpointScoringSetABoard1(t *testing.T) {
	// Sheet example: NS scores +120/+110/+50/+50/−100/−300 → MPs 10/8/5/5/2/0
	results := []BoardTraveller{
		{NsTeamId: "3", EwTeamId: "4", NsPoints: 120},
		{NsTeamId: "5", EwTeamId: "8", NsPoints: 110},
		{NsTeamId: "7", EwTeamId: "2", NsPoints: 50},
		{NsTeamId: "6", EwTeamId: "10", NsPoints: 50},
		{NsTeamId: "11", EwTeamId: "9", NsPoints: -100},
		{NsTeamId: "12", EwTeamId: "1", NsPoints: -300},
	}

	awards := MatchpointScoring{}.ScoreBoard(results)
	expected := []struct {
		ns, ew float64
	}{
		{10, 0},
		{8, 2},
		{5, 5},
		{5, 5},
		{2, 8},
		{0, 10},
	}

	if len(awards) != len(expected) {
		t.Fatalf("expected %d awards, got %d", len(expected), len(awards))
	}
	for i, exp := range expected {
		if awards[i].NsAward != exp.ns || awards[i].EwAward != exp.ew {
			t.Fatalf("row %d: expected NS/EW %.0f/%.0f, got %.0f/%.0f", i, exp.ns, exp.ew, awards[i].NsAward, awards[i].EwAward)
		}
	}
}

func TestCalculateTournamentScoreDenseRank(t *testing.T) {
	results := []PlayedResult{
		{DealNo: 1, NsTeamId: "a", EwTeamId: "b", Contract: "2N", Tricks: 8, Declarer: "N"},
		{DealNo: 1, NsTeamId: "c", EwTeamId: "d", Contract: "2H", Tricks: 8, Declarer: "N"},
	}

	score, err := CalculateTournamentScore("t1", MatchpointScoring{}, results)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(score.State.Standings) != 4 {
		t.Fatalf("expected 4 standings, got %d", len(score.State.Standings))
	}
	if score.State.Standings[0].Rank != 1 || score.State.Standings[0].TeamId != "a" {
		t.Fatalf("unexpected top standing: %+v", score.State.Standings[0])
	}
	if score.State.Standings[0].Percentage != 100 {
		t.Fatalf("expected 100%% for top score, got %v", score.State.Standings[0].Percentage)
	}
}
