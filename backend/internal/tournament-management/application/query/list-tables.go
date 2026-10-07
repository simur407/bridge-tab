package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
)

type ListTablesQuery struct {
	TournamentId string
}

func (c *ListTablesQuery) Execute(repo domain.TableReadRepository) ([]domain.TableDto, error) {
	return repo.FindAll(&c.TournamentId)
}
