package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
)

type AddScheduleCommand struct {
	TournamentId string
	RoundNo      int
	TableId      string
	NsTeamId     string
	EwTeamId     string
	SetId        string
}

func (c *AddScheduleCommand) Execute(repo domain.TournamentRepository) error {
	tournamentId := domain.TournamentId(c.TournamentId)
	tableId := domain.TableId(c.TableId)
	nsTeamId := domain.TeamId(c.NsTeamId)
	ewTeamId := domain.TeamId(c.EwTeamId)
	setId := domain.SetId(c.SetId)

	t, err := repo.Load(&tournamentId)
	if err != nil {
		return err
	}

	if err := t.AddSchedule(c.RoundNo, tableId, nsTeamId, ewTeamId, setId); err != nil {
		return err
	}

	return repo.Save(t)
}
