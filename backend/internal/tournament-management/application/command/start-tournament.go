package tournament_management

import (
	play "bridge-tab/internal/play/application"
	play_domain "bridge-tab/internal/play/domain"
	query "bridge-tab/internal/tournament-management/application/query"
	domain "bridge-tab/internal/tournament-management/domain"
)

type StartTurnamentCommand struct {
	TournamentId string
}

// Execute executes the command
func (c *StartTurnamentCommand) Execute(
	tournamentRepo domain.TournamentRepository,
	teamReadRepo domain.TeamReadRepository,
	boardProtocolReadRepo domain.BoardProtocolReadRepository,
	playSessionRepo play_domain.PlaySessionRepository,
) error {
	id := domain.TournamentId(c.TournamentId)
	// Load the turnament from the repository
	t, err := tournamentRepo.Load(&id)
	if err != nil {
		return err
	}

	if err := t.Start(); err != nil {
		return err
	}

	// Save the turnament to the repository
	err = tournamentRepo.Save(t)
	if err != nil {
		return err
	}

	listTeams := query.ListTeamsQuery{TournamentId: c.TournamentId}
	teams, err := listTeams.Execute(teamReadRepo)
	if err != nil {
		return err
	}

	listBoardProtocols := query.ListBoardProtocolsQuery{TournamentId: c.TournamentId}
	boardProtocols, err := listBoardProtocols.Execute(boardProtocolReadRepo)
	if err != nil {
		return err
	}

	playTeams := make([]play.TeamDto, len(teams))
	for i, team := range teams {
		players := make([]string, len(team.Members))
		for j, member := range team.Members {
			players[j] = string(member.Id)
		}
		playTeams[i] = play.TeamDto{Id: team.Id, Number: team.Number, Players: players}
	}

	matches := []play.MatchDto{}
	for _, boardProtocol := range boardProtocols {
		for _, teamPairs := range boardProtocol.TeamPairs {
			matches = append(matches, play.MatchDto{
				DealNo:     boardProtocol.BoardNo,
				NsTeam:     teamPairs.NS,
				EwTeam:     teamPairs.EW,
				Vulnerable: boardProtocol.Vulnerable,
			})
		}
	}

	provisionPlaySession := play.ProvisionPlaySessionCommand{
		TournamentId: c.TournamentId,
		Teams:        playTeams,
		Matches:      matches,
	}
	return provisionPlaySession.Execute(playSessionRepo)
}
