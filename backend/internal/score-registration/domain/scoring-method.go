package score_registration

type ScoringMethodName string

const (
	MethodMatchpoints ScoringMethodName = "matchpoints"
	MethodIMPs        ScoringMethodName = "imps"
)

// BoardTraveller is one comparison line on a deal (NS perspective points already computed).
type BoardTraveller struct {
	NsTeamId string
	EwTeamId string
	NsPoints int
}

// BoardAward is the scoring-method award for one traveller line.
type BoardAward struct {
	NsTeamId string
	EwTeamId string
	NsPoints int
	NsAward  float64
	EwAward  float64
}

// ScoringMethod converts raw NS bridge points on a board into awards (MP, IMP, …).
type ScoringMethod interface {
	Name() ScoringMethodName
	ScoreBoard(results []BoardTraveller) []BoardAward
}

func NewScoringMethod(name ScoringMethodName) ScoringMethod {
	switch name {
	case MethodIMPs:
		return ImpScoring{}
	default:
		return MatchpointScoring{}
	}
}
