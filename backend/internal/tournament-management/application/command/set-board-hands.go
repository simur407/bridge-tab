package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
)

type SetBoardHandsCommand struct {
	TournamentId string
	BoardNo      int
	North        string
	East         string
	South        string
	West         string
	Overwrite    bool
}

func (c *SetBoardHandsCommand) Execute(repo domain.TournamentRepository) error {
	tournamentId := domain.TournamentId(c.TournamentId)
	t, err := repo.Load(&tournamentId)
	if err != nil {
		return err
	}

	if err := t.SetBoardHands(c.BoardNo, c.North, c.East, c.South, c.West, c.Overwrite); err != nil {
		return err
	}

	return repo.Save(t)
}
