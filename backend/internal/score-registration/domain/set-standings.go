package score_registration

import "sort"

// TeamAwardStanding is a ranking built from board awards, using the same
// competition ranking as CalculateTournamentScore.
type TeamAwardStanding struct {
	TeamId     string
	TotalScore float64
	Percentage float64
	Rank       int
}

// RankBoardScores sums NS and EW awards per team and ranks them.
// The maximum for a deal is taken once, from the first line of that deal.
func RankBoardScores(scores []BoardScoreDto) []TeamAwardStanding {
	if len(scores) == 0 {
		return nil
	}

	totals := map[string]float64{}
	seenDeal := map[int]bool{}
	var maxPossible float64
	for _, score := range scores {
		totals[score.NsTeamId] += score.NsAward
		totals[score.EwTeamId] += score.EwAward
		if !seenDeal[score.DealNo] {
			seenDeal[score.DealNo] = true
			maxPossible += score.NsAward + score.EwAward
		}
	}

	standings := make([]TeamAwardStanding, 0, len(totals))
	for teamId, total := range totals {
		percentage := 0.0
		if maxPossible > 0 {
			percentage = 100.0 * total / maxPossible
		}
		standings = append(standings, TeamAwardStanding{
			TeamId:     teamId,
			TotalScore: total,
			Percentage: percentage,
		})
	}

	sort.Slice(standings, func(i, j int) bool {
		if standings[i].TotalScore == standings[j].TotalScore {
			return standings[i].TeamId < standings[j].TeamId
		}
		return standings[i].TotalScore > standings[j].TotalScore
	})

	for i := range standings {
		if i > 0 && standings[i].TotalScore == standings[i-1].TotalScore {
			standings[i].Rank = standings[i-1].Rank
		} else {
			standings[i].Rank = i + 1
		}
	}

	return standings
}
