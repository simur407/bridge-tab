package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
)

type GetTournamentById struct {
	Id string
}

func (c *GetTournamentById) Execute(repo domain.TournamentReadRepository) (*domain.TournamentDto, error) {
	t, err := repo.FindById(c.Id)
	if err != nil {
		return nil, err
	}

	return t, nil
}
