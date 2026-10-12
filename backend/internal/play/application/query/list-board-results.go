package play

import (
	domain "bridge-tab/internal/play/domain"
)

type ListBoardResultsQuery struct {
	TournamentId string
}

func (query *ListBoardResultsQuery) Execute(repository domain.BoardResultReadRepository) ([]domain.BoardResultDto, error) {
	return repository.FindAll(query.TournamentId)
}
