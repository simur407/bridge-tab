package rounds_registration

import (
	domain "bridge-tab/internal/rounds-registration/domain"
	"fmt"

	"context"
	"database/sql"
	"errors"
	"slices"
)

var ErrMisconfiguredGameSession = errors.New("game session is not configured correctly")
var ErrGameSessionNotFound = errors.New("game session not found")

type PostgresGameSessionRepository struct {
	Ctx context.Context
	Tx  *sql.Tx
}

func (r *PostgresGameSessionRepository) Load(id *domain.GameSessionId) (*domain.GameSession, error) {
	var gameSession domain.GameSession

	rows, err := r.Tx.QueryContext(r.Ctx, `
		SELECT 
			 team_id,
			 player_id
		FROM rounds_registration.team_players 
		WHERE game_session_id = $1`, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var teamId domain.TeamId
		var playerId domain.PlayerId
		err := rows.Scan(&teamId, &playerId)

		if err != nil {
			return nil, err
		}

		teamIndex := slices.IndexFunc(gameSession.State.Teams, func(tt *domain.Team) bool {
			return tt.Id == teamId
		})
		if teamIndex == -1 {
			gameSession.State.Teams = append(gameSession.State.Teams, &domain.Team{Id: teamId, Players: []domain.PlayerId{playerId}})
		} else {
			team := gameSession.State.Teams[teamIndex]
			team.Players = append(team.Players, playerId)
		}
	}

	rows, err = r.Tx.QueryContext(r.Ctx, `
		SELECT 
			deal_no,
			ns_team_id,
			ew_team_id,
			contract,
			declarer,
			tricks,
			opening_lead 
		FROM rounds_registration.round 
		WHERE game_session_id = $1 
		ORDER BY deal_no ASC`, id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var round domain.Round
		var nsTeam, ewTeam domain.TeamId
		var contract, declarer, openingLead sql.NullString
		var tricks sql.NullInt16

		err := rows.Scan(&round.DealNo, &nsTeam, &ewTeam, &contract, &declarer, &tricks, &openingLead)

		nsTeamIndex := slices.IndexFunc(gameSession.State.Teams, func(tt *domain.Team) bool {
			return tt.Id == nsTeam
		})
		ewTeamIndex := slices.IndexFunc(gameSession.State.Teams, func(tt *domain.Team) bool {
			return tt.Id == ewTeam
		})

		if nsTeamIndex == -1 || ewTeamIndex == -1 {
			return nil, ErrMisconfiguredGameSession
		}

		round.NsTeam = &gameSession.State.Teams[nsTeamIndex].Id
		round.EwTeam = &gameSession.State.Teams[ewTeamIndex].Id

		if contract.Valid {
			round.Contract = contract.String
		}
		if declarer.Valid {
			round.Declarer = declarer.String
		}
		if tricks.Valid {
			round.Tricks = int(tricks.Int16)
		}
		if openingLead.Valid {
			round.OpeningLead = openingLead.String
		}
		if err != nil {
			return nil, err
		}
		gameSession.State.Rounds = append(gameSession.State.Rounds, &round)
	}

	if len(gameSession.State.Teams) == 0 && len(gameSession.State.Rounds) == 0 {
		return nil, ErrGameSessionNotFound
	}

	gameSession.State.Id = *id

	return &gameSession, nil
}

func (r *PostgresGameSessionRepository) Save(gameSession *domain.GameSession) error {
	for _, event := range gameSession.GetEvents() {
		switch event := event.(type) {
		case domain.GameSessionStarted:
			return r.gameSessionStarted(event)
		case domain.RoundPlayed:
			return r.roundPlayed(event)
		default:
			return errors.New("unknown event")
		}
	}

	return nil
}

func (r *PostgresGameSessionRepository) gameSessionStarted(event domain.GameSessionStarted) error {
	var teamsQuery string
	var teams []interface{}

	var paramDepth = 0
	for _, team := range event.Teams {
		for _, player := range team.Players {
			teams = append(teams, event.GameSessionId, team.Id, player)
			teamsQuery += fmt.Sprintf("($%d, $%d, $%d), ", paramDepth*3+1, paramDepth*3+2, paramDepth*3+3)
			paramDepth++
		}
	}
	teamsQuery = teamsQuery[:len(teamsQuery)-2]

	_, err := r.Tx.ExecContext(r.Ctx, `
		INSERT INTO rounds_registration.team_players (game_session_id, team_id, player_id) VALUES `+teamsQuery, teams...)

	if err != nil {
		return err
	}

	var roundsQuery string
	var rounds []interface{}
	for i, round := range event.Rounds {
		rounds = append(rounds, event.GameSessionId, round.DealNo, round.NsTeam, round.EwTeam)
		roundsQuery += fmt.Sprintf("($%d, $%d, $%d, $%d), ", i*4+1, i*4+2, i*4+3, i*4+4)
	}
	roundsQuery = roundsQuery[:len(roundsQuery)-2]

	_, err = r.Tx.ExecContext(r.Ctx, `INSERT INTO rounds_registration.round (game_session_id, deal_no, ns_team_id, ew_team_id) VALUES `+roundsQuery, rounds...)

	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresGameSessionRepository) roundPlayed(event domain.RoundPlayed) error {
	_, err := r.Tx.ExecContext(r.Ctx, `
		UPDATE rounds_registration.round 
		SET contract = $5, tricks = $6, declarer = $7, opening_lead = $8, updated_at = now()
		WHERE game_session_id = $1 AND deal_no = $2 AND ns_team_id = $3 AND ew_team_id = $4`,
		event.GameSessionId,
		event.DealNo,
		event.NsTeamId,
		event.EwTeamId,
		event.Contract,
		event.Tricks,
		event.Declarer,
		event.OpeningLead,
	)

	if err != nil {
		return err
	}

	return nil
}
