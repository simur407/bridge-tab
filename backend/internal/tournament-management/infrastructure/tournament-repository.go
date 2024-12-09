package tournament_management

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"

	domain "bridge-tab/internal/tournament-management/domain"
)

type PostgresTournamentRepository struct {
	Ctx context.Context
	Tx  *sql.Tx
}

var ErrTournamentNotFound = errors.New("tournament not found")

func (r *PostgresTournamentRepository) Load(Id *domain.TournamentId) (*domain.Tournament, error) {
	var Tournament domain.Tournament
	row := r.Tx.QueryRowContext(r.Ctx, "SELECT id, name FROM tournament_management.tournament WHERE id = $1", Id)
	err := row.Scan(&Tournament.State.Id, &Tournament.State.Name)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: %v", ErrTournamentNotFound, Id)
		}
		return nil, err
	}

	contestantRows, err := r.Tx.QueryContext(r.Ctx, "SELECT id FROM tournament_management.contestant WHERE tournament_id = $1", Id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	var Contestants []*domain.Contestant
	for contestantRows.Next() {
		var contestant domain.Contestant
		err = contestantRows.Scan(&contestant.Id)
		if err != nil {
			return nil, err
		}
		Contestants = append(Contestants, &contestant)
	}

	teamRows, err := r.Tx.QueryContext(r.Ctx, "SELECT id, name FROM tournament_management.team WHERE Tournament_id = $1", Id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	var Teams []*domain.Team
	for teamRows.Next() {
		var team domain.Team
		err = teamRows.Scan(&team.State.Id, &team.State.Name)
		if err != nil {
			return nil, err
		}
		team.State.TournamentId = *Id
		Teams = append(Teams, &team)
	}

	teamContestantRows, err := r.Tx.QueryContext(r.Ctx, `
	SELECT team_id, contestant_id FROM tournament_management.team_contestant 
	INNER JOIN tournament_management.team ON team_contestant.team_id = team.id 
	WHERE team.tournament_id = $1`, Id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	for teamContestantRows.Next() {
		var teamContestant struct {
			ContestantId string
			TeamId       string
		}
		err = teamContestantRows.Scan(&teamContestant.TeamId, &teamContestant.ContestantId)
		if err != nil {
			return nil, err
		}
		for _, team := range Teams {
			if team.State.Id == domain.TeamId(teamContestant.TeamId) {
				contestantIndex := slices.IndexFunc(Contestants, func(c *domain.Contestant) bool {
					return c.Id == domain.ContestantId(teamContestant.ContestantId)
				})

				if contestantIndex != -1 {
					team.State.Members = append(team.State.Members, Contestants[contestantIndex])
					Contestants[contestantIndex].Team = team
				}
			}
		}
	}

	Tournament.State.Contestants = Contestants
	Tournament.State.Teams = Teams

	return &Tournament, nil
}

func (r *PostgresTournamentRepository) Save(t *domain.Tournament) error {
	for _, event := range t.GetEvents() {
		switch event := event.(type) {
		case domain.TournamentCreated:
			return r.TournamentCreated(event)
		case domain.TournamentRemoved:
			return r.TournamentRemoved(event)
		case domain.TournamentStarted:
			return r.TournamentStarted(event)
		case domain.ContestantJoinedTournament:
			return r.contestantJoinedTournament(event)
		case domain.ContestantLeftTournament:
			return r.contestantLeftTournament(event)
		case domain.TeamCreated:
			return r.teamCreated(event)
		case domain.TeamRemoved:
			return r.teamRemoved(event)
		case domain.ContestantJoinedTeam:
			return r.contestantJoinedTeam(event)
		case domain.ContestantLeftTeam:
			return r.contestantLeftTeam(event)

		default:
			return errors.New("unknown event")
		}
	}
	t.Commit()
	return nil
}

func (r *PostgresTournamentRepository) TournamentCreated(event domain.TournamentCreated) error {
	_, err := r.Tx.ExecContext(r.Ctx, "INSERT INTO tournament_management.tournament (id, name) VALUES ($1, $2)", event.TournamentId, event.Name)

	return err
}

func (r *PostgresTournamentRepository) TournamentRemoved(event domain.TournamentRemoved) error {
	_, err := r.Tx.ExecContext(r.Ctx, "DELETE FROM tournament_management.tournament WHERE id = $1", event.TournamentId)

	return err
}

func (r *PostgresTournamentRepository) TournamentStarted(event domain.TournamentStarted) error {
	_, err := r.Tx.ExecContext(r.Ctx, "UPDATE tournament_management.tournament SET started_at = $1 WHERE id = $2", event.StartedAt, event.TournamentId)

	return err
}

func (r *PostgresTournamentRepository) contestantJoinedTournament(event domain.ContestantJoinedTournament) error {
	_, err := r.Tx.ExecContext(r.Ctx, "INSERT INTO tournament_management.contestant (id, tournament_id) VALUES ($1, $2)", event.ContestantId, event.TournamentId)

	return err
}

func (r *PostgresTournamentRepository) contestantLeftTournament(event domain.ContestantLeftTournament) error {
	_, err := r.Tx.ExecContext(r.Ctx, "DELETE FROM tournament_management.contestant WHERE id = $1 AND tournament_id = $2", event.ContestantId, event.TournamentId)

	return err
}

func (r *PostgresTournamentRepository) teamCreated(event domain.TeamCreated) error {
	_, err := r.Tx.ExecContext(r.Ctx, "INSERT INTO tournament_management.team (id, Tournament_id, name) VALUES ($1, $2, $3)", event.TeamId, event.TournamentId, event.Name)

	return err
}

func (r *PostgresTournamentRepository) teamRemoved(event domain.TeamRemoved) error {
	_, err := r.Tx.ExecContext(r.Ctx, "DELETE FROM tournament_management.team WHERE id = $1 AND tournament_id = $2", event.TeamId, event.TournamentId)

	return err
}

func (r *PostgresTournamentRepository) contestantJoinedTeam(event domain.ContestantJoinedTeam) error {
	_, err := r.Tx.ExecContext(r.Ctx, "INSERT INTO tournament_management.team_contestant (team_id, contestant_id) VALUES ($1, $2)", event.TeamId, event.ContestantId)

	return err
}

func (r *PostgresTournamentRepository) contestantLeftTeam(event domain.ContestantLeftTeam) error {
	_, err := r.Tx.ExecContext(r.Ctx, "DELETE FROM tournament_management.team_contestant WHERE team_id = $1 AND contestant_id = $2", event.TeamId, event.ContestantId)

	return err
}
