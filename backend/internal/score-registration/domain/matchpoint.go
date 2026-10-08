package score_registration

// MatchpointScoring implements classic duplicate matchpoints (2 for a win, 1 for a tie).
type MatchpointScoring struct{}

func (MatchpointScoring) Name() ScoringMethodName {
	return MethodMatchpoints
}

func (MatchpointScoring) ScoreBoard(results []BoardTraveller) []BoardAward {
	n := len(results)
	awards := make([]BoardAward, n)
	if n == 0 {
		return awards
	}

	for i, result := range results {
		worse := 0
		equal := 0
		for j, other := range results {
			if i == j {
				continue
			}
			if other.NsPoints < result.NsPoints {
				worse++
			} else if other.NsPoints == result.NsPoints {
				equal++
			}
		}
		nsMP := float64(2*worse + equal)
		top := float64(2 * (n - 1))
		awards[i] = BoardAward{
			NsTeamId: result.NsTeamId,
			EwTeamId: result.EwTeamId,
			NsPoints: result.NsPoints,
			NsAward:  nsMP,
			EwAward:  top - nsMP,
		}
	}
	return awards
}
