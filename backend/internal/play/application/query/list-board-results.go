package play

import (
	domain "bridge-tab/internal/play/domain"
)

type ListBoardResultsQuery struct {
	PlaySessionId string
}

func (query *ListBoardResultsQuery) Execute(repository domain.BoardResultReadRepository) ([]domain.BoardResultDto, error) {
	playSessionId := query.PlaySessionId
	return repository.FindAll(playSessionId)
}
