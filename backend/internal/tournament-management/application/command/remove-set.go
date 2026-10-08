package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
)

type RemoveSetCommand struct {
	TournamentId string
	SetId        string
}

func (c *RemoveSetCommand) Execute(repo domain.TournamentRepository) error {
	id := domain.TournamentId(c.TournamentId)
	t, err := repo.Load(&id)
	if err != nil {
		return err
	}

	if err := t.RemoveSet(domain.SetId(c.SetId)); err != nil {
		return err
	}

	return repo.Save(t)
}
