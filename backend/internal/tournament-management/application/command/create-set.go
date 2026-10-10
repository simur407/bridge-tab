package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
)

type CreateSetCommand struct {
	TournamentId string
	SetId        string
	Label        string
	BoardNos     []int
}

func (c *CreateSetCommand) Execute(repo domain.TournamentRepository) error {
	id := domain.TournamentId(c.TournamentId)
	t, err := repo.Load(&id)
	if err != nil {
		return err
	}

	if err := t.CreateSet(domain.SetId(c.SetId), c.Label, c.BoardNos); err != nil {
		return err
	}

	return repo.Save(t)
}
