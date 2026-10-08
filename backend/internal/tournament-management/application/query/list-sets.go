package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
)

type ListSetsQuery struct {
	TournamentId string
}

func (q *ListSetsQuery) Execute(repo domain.SetReadRepository) ([]domain.SetDto, error) {
	return repo.FindAll(&q.TournamentId)
}
