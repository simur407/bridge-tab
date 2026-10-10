package score_registration

// ImpScoring is a placeholder seam for future IMP scoring.
// It currently returns zero awards so the method can be selected without panicking.
type ImpScoring struct{}

func (ImpScoring) Name() ScoringMethodName {
	return MethodIMPs
}

func (ImpScoring) ScoreBoard(results []BoardTraveller) []BoardAward {
	awards := make([]BoardAward, len(results))
	for i, result := range results {
		awards[i] = BoardAward{
			NsTeamId: result.NsTeamId,
			EwTeamId: result.EwTeamId,
			NsPoints: result.NsPoints,
		}
	}
	return awards
}
