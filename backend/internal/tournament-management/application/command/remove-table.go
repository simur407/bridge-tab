package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
)

type RemoveTableCommand struct {
	TournamentId string
	TableId      string
}

func (c *RemoveTableCommand) Execute(repo domain.TournamentRepository) error {
	id := domain.TournamentId(c.TournamentId)
	tableId := domain.TableId(c.TableId)
	t, err := repo.Load(&id)
	if err != nil {
		return err
	}

	if err := t.DeleteTable(&tableId); err != nil {
		return err
	}

	return repo.Save(t)
}
