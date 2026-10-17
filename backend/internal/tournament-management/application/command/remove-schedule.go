package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
)

type RemoveScheduleCommand struct {
	TournamentId string
	RoundNo      int
	TableId      string
}

func (c *RemoveScheduleCommand) Execute(repo domain.TournamentRepository) error {
	tournamentId := domain.TournamentId(c.TournamentId)
	tableId := domain.TableId(c.TableId)

	t, err := repo.Load(&tournamentId)
	if err != nil {
		return err
	}

	if err := t.RemoveSchedule(c.RoundNo, tableId); err != nil {
		return err
	}

	return repo.Save(t)
}
