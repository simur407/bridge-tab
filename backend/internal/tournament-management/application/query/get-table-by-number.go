package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
)

type GetTableByNumberQuery struct {
	TournamentId string
	Number       int
}

func (q *GetTableByNumberQuery) Execute(repo domain.TableReadRepository) (*domain.TableDto, error) {
	return repo.FindByNumber(&q.TournamentId, q.Number)
}
