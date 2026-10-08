package score_registration

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	domain "bridge-tab/internal/score-registration/domain"
)

type PostgresTournamentScoreRepository struct {
	Ctx context.Context
	Tx  *sql.Tx
}

func (r *PostgresTournamentScoreRepository) Load(id *domain.TournamentScoreId, method domain.ScoringMethodName) (*domain.TournamentScore, error) {
	row := r.Tx.QueryRowContext(r.Ctx, `
		SELECT tournament_id, method
		FROM score_registration.tournament_score
		WHERE tournament_id = $1 AND method = $2`, id, method)

	var score domain.TournamentScore
	var storedMethod string
	err := row.Scan(&score.State.Id, &storedMethod)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("tournament score not found: %v (%s)", id, method)
		}
		return nil, err
	}
	score.State.Method = domain.ScoringMethodName(storedMethod)

	standingRows, err := r.Tx.QueryContext(r.Ctx, `
		SELECT team_id, total_score, percentage, rank
		FROM score_registration.standing
		WHERE tournament_id = $1 AND method = $2
		ORDER BY rank, team_id`, id, method)
	if err != nil {
		return nil, err
	}
	for standingRows.Next() {
		var standing domain.Standing
		if err := standingRows.Scan(&standing.TeamId, &standing.TotalScore, &standing.Percentage, &standing.Rank); err != nil {
			return nil, err
		}
		score.State.Standings = append(score.State.Standings, standing)
	}

	boardRows, err := r.Tx.QueryContext(r.Ctx, `
		SELECT deal_no, ns_team_id, ew_team_id, ns_points, ns_award, ew_award
		FROM score_registration.board_score
		WHERE tournament_id = $1 AND method = $2
		ORDER BY deal_no, ns_team_id`, id, method)
	if err != nil {
		return nil, err
	}
	for boardRows.Next() {
		var board domain.BoardScore
		if err := boardRows.Scan(&board.DealNo, &board.NsTeamId, &board.EwTeamId, &board.NsPoints, &board.NsAward, &board.EwAward); err != nil {
			return nil, err
		}
		score.State.BoardScores = append(score.State.BoardScores, board)
	}

	return &score, nil
}

func (r *PostgresTournamentScoreRepository) Save(score *domain.TournamentScore) error {
	for _, event := range score.GetEvents() {
		switch event := event.(type) {
		case domain.TournamentScoreCalculated:
			if err := r.tournamentScoreCalculated(event); err != nil {
				return err
			}
		default:
			return errors.New("unknown event")
		}
	}
	score.Commit()
	return nil
}

func (r *PostgresTournamentScoreRepository) tournamentScoreCalculated(event domain.TournamentScoreCalculated) error {
	_, err := r.Tx.ExecContext(r.Ctx, `
		DELETE FROM score_registration.tournament_score
		WHERE tournament_id = $1 AND method = $2`, event.TournamentId, event.Method)
	if err != nil {
		return err
	}

	_, err = r.Tx.ExecContext(r.Ctx, `
		INSERT INTO score_registration.tournament_score (tournament_id, method)
		VALUES ($1, $2)`, event.TournamentId, event.Method)
	if err != nil {
		return err
	}

	for _, standing := range event.Standings {
		_, err = r.Tx.ExecContext(r.Ctx, `
			INSERT INTO score_registration.standing (tournament_id, method, team_id, total_score, percentage, rank)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			event.TournamentId, event.Method, standing.TeamId, standing.TotalScore, standing.Percentage, standing.Rank)
		if err != nil {
			return err
		}
	}

	for _, board := range event.BoardScores {
		_, err = r.Tx.ExecContext(r.Ctx, `
			INSERT INTO score_registration.board_score (
				tournament_id, method, deal_no, ns_team_id, ew_team_id, ns_points, ns_award, ew_award
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			event.TournamentId, event.Method, board.DealNo, board.NsTeamId, board.EwTeamId,
			board.NsPoints, board.NsAward, board.EwAward)
		if err != nil {
			return err
		}
	}

	return nil
}
