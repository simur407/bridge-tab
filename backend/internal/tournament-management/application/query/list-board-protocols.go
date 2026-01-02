package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
)

type ListBoardProtocolsQuery struct {
	TournamentId string
}

func (q *ListBoardProtocolsQuery) Execute(repo domain.BoardProtocolReadRepository) ([]domain.BoardProtocolDto, error) {
	bp, err := repo.FindAll(&q.TournamentId)
	if err != nil {
		return nil, err
	}

	return bp, nil
}
