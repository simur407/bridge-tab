package score_registration

import (
	domain "bridge-tab/internal/score-registration/domain"
)

type GetStandingsQuery struct {
	TournamentId string
	Method       domain.ScoringMethodName
}

func (q *GetStandingsQuery) Execute(repo domain.TournamentScoreReadRepository) ([]domain.StandingDto, error) {
	return repo.FindStandings(q.TournamentId, q.Method)
}
