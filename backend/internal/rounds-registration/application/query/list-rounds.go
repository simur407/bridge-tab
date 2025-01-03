package rounds_registration

import (
	domain "bridge-tab/internal/rounds-registration/domain"
)

type ListRoundsQuery struct {
	GameSessionId string
}

func (query *ListRoundsQuery) Execute(repository domain.GameSessionReadRepository) ([]domain.PlayedRoundDto, error) {
	rounds, err := repository.FindAllRounds(&query.GameSessionId)

	if err != nil {
		return nil, err
	}

	return rounds, nil
}
