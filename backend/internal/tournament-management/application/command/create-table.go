package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
)

type CreateTableCommand struct {
	TournamentId string
	TableId      string
	Number       int
}

func (c *CreateTableCommand) Execute(repo domain.TournamentRepository) error {
	tournamentId := domain.TournamentId(c.TournamentId)
	tableId := domain.TableId(c.TableId)
	t, err := repo.Load(&tournamentId)
	if err != nil {
		return err
	}

	if err := t.CreateTable(&tableId, c.Number); err != nil {
		return err
	}

	return repo.Save(t)
}
