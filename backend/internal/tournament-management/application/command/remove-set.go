package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
)

type RemoveSetCommand struct {
	TournamentId string
	SetId        string
}

func (c *RemoveSetCommand) Execute(repo domain.TournamentRepository) error {
	tournamentId := domain.TournamentId(c.TournamentId)
	setId := domain.SetId(c.SetId)

	t, err := repo.Load(&tournamentId)
	if err != nil {
		return err
	}

	if err := t.RemoveSet(setId); err != nil {
		return err
	}

	return repo.Save(t)
}
