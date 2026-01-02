package tournament_management

import (
	round_registration "bridge-tab/internal/rounds-registration/application"
	round_registration_domain "bridge-tab/internal/rounds-registration/domain"
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
	gameSessionRepo round_registration_domain.GameSessionRepository) error {
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

	sessionTeams := make([]round_registration.TeamDto, len(teams))
	for i, team := range teams {
		players := make([]string, len(team.Members))
		for j, member := range team.Members {
			players[j] = string(member.Id)
		}
		sessionTeams[i] = round_registration.TeamDto{Id: team.Id, Players: players}
	}

	rounds := []round_registration.RoundDto{}
	for _, boardProtocol := range boardProtocols {
		for _, teamPairs := range boardProtocol.TeamPairs {
			rounds = append(rounds, round_registration.RoundDto{DealNo: boardProtocol.BoardNo, NsTeam: teamPairs.NS, EwTeam: teamPairs.EW})
		}
	}

	// Start the game session
	startGameSession := round_registration.StartGameSessionCommand{
		Id:     c.TournamentId,
		Teams:  sessionTeams,
		Rounds: rounds,
	}

	err = startGameSession.Execute(gameSessionRepo)
	if err != nil {
		return err
	}

	return nil
}
