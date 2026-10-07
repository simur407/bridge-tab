package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
)

type GetTableByIdQuery struct {
	TournamentId string
	TableId      string
}

func (q *GetTableByIdQuery) Execute(repo domain.TableReadRepository) (*domain.TableDto, error) {
	return repo.FindById(&q.TournamentId, &q.TableId)
}
