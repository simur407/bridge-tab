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
	tournamentId := domain.TournamentId(c.TournamentId)
	setId := domain.SetId(c.SetId)

	t, err := repo.Load(&tournamentId)
	if err != nil {
		return err
	}

	if err := t.CreateSet(setId, c.Label, c.BoardNos); err != nil {
		return err
	}

	return repo.Save(t)
}
