package e2e_test

import (
	"fmt"
	"testing"

	"bridge-tab/e2e/flow"
	rounds_registration_domain "bridge-tab/internal/rounds-registration/domain"
	score_registration_domain "bridge-tab/internal/score-registration/domain"
	tournament_domain "bridge-tab/internal/tournament-management/domain"
)

func TestTournamentHappyPath(t *testing.T) {
	flow := flow.New(t)

	flow.CreateTournament()
	flow.CreateTeam(1)
	flow.CreateTeam(2)
	flow.CreateTable(1)

	flow.CreateBoardProtocol(1, tournament_domain.None, 1, 2)
	flow.CreateBoardProtocolWithTable(2, tournament_domain.NS, 1, 1, 2)
	flow.CreateBoardProtocol(3, tournament_domain.EW, 1, 2)
	flow.CreateBoardProtocol(4, tournament_domain.Both, 1, 2)

	alice := flow.RegisterUser("Alice")
	flow.JoinTournament(alice)
	flow.JoinTeam(alice, 1)

	bob := flow.RegisterUser("Bob")
	flow.JoinTournament(bob)
	flow.JoinTeam(bob, 2)

	flow.StartTournament()
	flow.AssertPendingRounds(4)

	flow.PlayRound(alice, 2, 1, "3NT", 9, "N", "KH")
	flow.PlayRound(bob, 1, 2, "4H", 10, "E", "AS")
	flow.PlayRound(alice, 2, 3, "1Cx", 7, "S", "QD")
	flow.PlayRound(bob, 1, 4, "Pass", 0, "", "")

	rounds := flow.ListRounds()
	assertDeal(t, rounds, 1, 1, 2, "3N", 9, "N", "KH")
	assertDeal(t, rounds, 2, 1, 2, "4H", 10, "E", "AS")
	assertDeal(t, rounds, 3, 1, 2, "1Cx", 7, "S", "QD")
	assertDeal(t, rounds, 4, 1, 2, "Pass", 0, "", "")
}

func TestTournamentScoring(t *testing.T) {
	flow := flow.New(t)

	flow.CreateTournament()
	for teamNo := 1; teamNo <= 4; teamNo++ {
		flow.CreateTeam(teamNo)
	}

	// Two tables play the same board (classic MP field of size 2).
	flow.CreateBoardProtocolWithPairs(1, tournament_domain.None, [][2]int{{1, 2}, {3, 4}})
	flow.CreateSet("A", []int{1})

	players := map[int]string{}
	for teamNo := 1; teamNo <= 4; teamNo++ {
		userId := flow.RegisterUser(fmt.Sprintf("Player%d", teamNo))
		flow.JoinTournament(userId)
		flow.JoinTeam(userId, teamNo)
		players[teamNo] = userId
	}

	flow.StartTournament()
	flow.AssertPendingRounds(2)

	// Table 1: 3NT making = +400 NS; Table 2: 3NT down 1 = -50 NS
	flow.PlayRound(players[1], 2, 1, "3NT", 9, "N", "KH")
	flow.PlayRound(players[3], 4, 1, "3NT", 8, "N", "AS")

	flow.FinishTournament()
	flow.CalculateStandings(score_registration_domain.MethodMatchpoints)

	standings := flow.GetStandings(score_registration_domain.MethodMatchpoints)
	if len(standings) != 4 {
		t.Fatalf("expected 4 standings, got %d", len(standings))
	}

	byNumber := map[int]score_registration_domain.StandingDto{}
	for _, standing := range standings {
		byNumber[standing.TeamNumber] = standing
	}

	// Top on the board is 2 MP. Team 1 (NS +400) and team 4 (EW of -50) get 2.
	if byNumber[1].TotalScore != 2 || byNumber[4].TotalScore != 2 {
		t.Fatalf("expected teams 1 and 4 to have 2 MP, got 1=%v 4=%v", byNumber[1].TotalScore, byNumber[4].TotalScore)
	}
	if byNumber[2].TotalScore != 0 || byNumber[3].TotalScore != 0 {
		t.Fatalf("expected teams 2 and 3 to have 0 MP, got 2=%v 3=%v", byNumber[2].TotalScore, byNumber[3].TotalScore)
	}
	if byNumber[1].Rank != 1 || byNumber[4].Rank != 1 {
		t.Fatalf("expected dense rank 1 for teams 1 and 4, got 1=%d 4=%d", byNumber[1].Rank, byNumber[4].Rank)
	}

	boardScores := flow.GetBoardScores(score_registration_domain.MethodMatchpoints, 1)
	if len(boardScores) != 2 {
		t.Fatalf("expected 2 board scores, got %d", len(boardScores))
	}
}

func assertDeal(
	t *testing.T,
	rounds []rounds_registration_domain.PlayedRoundDto,
	dealNo, nsTeamNumber, ewTeamNumber int,
	contract string,
	tricks int,
	declarer, openingLead string,
) {
	t.Helper()

	var deal *rounds_registration_domain.PlayedRoundDto
	for i := range rounds {
		if rounds[i].DealNo == dealNo {
			deal = &rounds[i]
			break
		}
	}
	if deal == nil {
		t.Fatalf("missing deal %d in final rounds", dealNo)
	}
	if deal.NsTeamNumber != nsTeamNumber || deal.EwTeamNumber != ewTeamNumber {
		t.Fatalf("deal %d teams: got NS=%d EW=%d, want NS=%d EW=%d",
			dealNo, deal.NsTeamNumber, deal.EwTeamNumber, nsTeamNumber, ewTeamNumber)
	}
	if deal.Contract != contract || deal.Tricks != tricks || deal.Declarer != declarer || deal.OpeningLead != openingLead {
		t.Fatalf("deal %d score mismatch: %#v", dealNo, deal)
	}
}
