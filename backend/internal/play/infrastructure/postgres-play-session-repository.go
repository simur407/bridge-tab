package play

import (
	"context"
	"database/sql"
	"errors"
	"slices"

	"bridge-tab/internal/idutil"
	domain "bridge-tab/internal/play/domain"
)

type PostgresPlaySessionRepository struct {
	Ctx context.Context
	Tx  *sql.Tx
}

func (r *PostgresPlaySessionRepository) Load(id domain.TournamentId) (*domain.PlaySession, error) {
	var session domain.PlaySession

	teamRows, err := r.Tx.QueryContext(r.Ctx, `
		SELECT team_id, number
		FROM play.team
		WHERE tournament_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer teamRows.Close()

	for teamRows.Next() {
		var team domain.Team
		if err := teamRows.Scan(&team.Id, &team.Number); err != nil {
			teamRows.Close()
			return nil, err
		}
		session.State.Teams = append(session.State.Teams, &team)
	}
	if err := teamRows.Err(); err != nil {
		teamRows.Close()
		return nil, err
	}
	teamRows.Close()

	playerRows, err := r.Tx.QueryContext(r.Ctx, `
		SELECT team_id, player_id
		FROM play.team_player
		WHERE tournament_id = $1`, id)
	if err != nil {
		return nil, err
	}

	for playerRows.Next() {
		var teamId domain.TeamId
		var playerId domain.PlayerId
		if err := playerRows.Scan(&teamId, &playerId); err != nil {
			return nil, err
		}
		index := slices.IndexFunc(session.State.Teams, func(team *domain.Team) bool {
			return idutil.SameId(team.Id, teamId)
		})
		if index == -1 {
			playerRows.Close()
			return nil, errors.New("play is not configured correctly")
		}
		session.State.Teams[index].Players = append(session.State.Teams[index].Players, playerId)
	}
	if err := playerRows.Err(); err != nil {
		playerRows.Close()
		return nil, err
	}
	playerRows.Close()

	matchRows, err := r.Tx.QueryContext(r.Ctx, `
		SELECT deal_no, ns_team_id, ew_team_id, vulnerable
		FROM play.match
		WHERE tournament_id = $1
		ORDER BY deal_no`, id)
	if err != nil {
		return nil, err
	}
	defer matchRows.Close()

	for matchRows.Next() {
		var match domain.Match
		if err := matchRows.Scan(&match.DealNo, &match.NsTeamId, &match.EwTeamId, &match.Vulnerable); err != nil {
			return nil, err
		}
		session.State.Matches = append(session.State.Matches, &match)
	}
	if err := matchRows.Err(); err != nil {
		return nil, err
	}

	if len(session.State.Teams) == 0 && len(session.State.Matches) == 0 {
		return nil, domain.ErrPlaySessionNotProvisioned
	}

	session.State.TournamentId = id
	return &session, nil
}

func (r *PostgresPlaySessionRepository) Save(session *domain.PlaySession) error {
	for _, event := range session.GetEvents() {
		var err error
		switch event := event.(type) {
		case domain.PlaySessionProvisioned:
			err = r.playSessionProvisioned(event)
		default:
			return errors.New("unknown event")
		}
		if err != nil {
			return err
		}
	}
	session.Commit()
	return nil
}

func (r *PostgresPlaySessionRepository) playSessionProvisioned(event domain.PlaySessionProvisioned) error {
	for _, team := range event.Teams {
		_, err := r.Tx.ExecContext(r.Ctx, `
			INSERT INTO play.team (tournament_id, team_id, number)
			VALUES ($1, $2, $3)`,
			event.TournamentId, team.Id, team.Number,
		)
		if err != nil {
			return err
		}
		for _, player := range team.Players {
			_, err = r.Tx.ExecContext(r.Ctx, `
				INSERT INTO play.team_player (tournament_id, team_id, player_id)
				VALUES ($1, $2, $3)`,
				event.TournamentId, team.Id, player,
			)
			if err != nil {
				return err
			}
		}
	}

	for _, match := range event.Matches {
		_, err := r.Tx.ExecContext(r.Ctx, `
			INSERT INTO play.match (tournament_id, deal_no, ns_team_id, ew_team_id, vulnerable)
			VALUES ($1, $2, $3, $4, $5)`,
			event.TournamentId, match.DealNo, match.NsTeamId, match.EwTeamId, match.Vulnerable,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
