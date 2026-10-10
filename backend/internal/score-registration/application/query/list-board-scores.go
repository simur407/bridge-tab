package score_registration

import (
	domain "bridge-tab/internal/score-registration/domain"
)

type ListBoardScoresQuery struct {
	TournamentId string
	Method       domain.ScoringMethodName
	DealNo       *int
}

func (q *ListBoardScoresQuery) Execute(repo domain.TournamentScoreReadRepository) ([]domain.BoardScoreDto, error) {
	return repo.FindBoardScores(q.TournamentId, q.Method, q.DealNo)
}
