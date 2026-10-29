package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
)

type SetScoreEntryModeCommand struct {
	TournamentId   string
	ScoreEntryMode string
}

func (c *SetScoreEntryModeCommand) Execute(repo domain.TournamentRepository) error {
	tournamentId := domain.TournamentId(c.TournamentId)
	t, err := repo.Load(&tournamentId)
	if err != nil {
		return err
	}

	mode, err := domain.ParseScoreEntryMode(c.ScoreEntryMode)
	if err != nil {
		return err
	}
	if err := t.SetScoreEntryMode(mode); err != nil {
		return err
	}

	return repo.Save(t)
}
