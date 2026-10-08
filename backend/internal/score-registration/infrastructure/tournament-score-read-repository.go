package score_registration

import (
	"context"
	"database/sql"

	domain "bridge-tab/internal/score-registration/domain"
)

type PostgresTournamentScoreReadRepository struct {
	Ctx context.Context
	Tx  *sql.Tx
}

func (r *PostgresTournamentScoreReadRepository) FindStandings(tournamentId string, method domain.ScoringMethodName) ([]domain.StandingDto, error) {
	rows, err := r.Tx.QueryContext(r.Ctx, `
		SELECT standing.team_id, team.number, team.name, standing.total_score, standing.percentage, standing.rank
		FROM score_registration.standing
		LEFT JOIN tournament_management.team
			ON standing.team_id = team.id
		WHERE standing.tournament_id = $1 AND standing.method = $2
		ORDER BY standing.rank, team.number`, tournamentId, method)
	if err != nil {
		return nil, err
	}

	var standings []domain.StandingDto
	for rows.Next() {
		var standing domain.StandingDto
		var teamNumber sql.NullInt64
		var teamName sql.NullString
		if err := rows.Scan(&standing.TeamId, &teamNumber, &teamName, &standing.TotalScore, &standing.Percentage, &standing.Rank); err != nil {
			return nil, err
		}
		if teamNumber.Valid {
			standing.TeamNumber = int(teamNumber.Int64)
		}
		if teamName.Valid {
			standing.TeamName = teamName.String
		}
		standings = append(standings, standing)
	}
	return standings, nil
}

func (r *PostgresTournamentScoreReadRepository) FindBoardScores(tournamentId string, method domain.ScoringMethodName, dealNo *int) ([]domain.BoardScoreDto, error) {
	query := `
		SELECT board_score.deal_no,
			board_score.ns_team_id, ns_team.number,
			board_score.ew_team_id, ew_team.number,
			board_score.ns_points, board_score.ns_award, board_score.ew_award
		FROM score_registration.board_score
		LEFT JOIN tournament_management.team AS ns_team ON board_score.ns_team_id = ns_team.id
		LEFT JOIN tournament_management.team AS ew_team ON board_score.ew_team_id = ew_team.id
		WHERE board_score.tournament_id = $1 AND board_score.method = $2`
	args := []any{tournamentId, method}
	if dealNo != nil {
		query += ` AND board_score.deal_no = $3`
		args = append(args, *dealNo)
	}
	query += ` ORDER BY board_score.deal_no, board_score.ns_award DESC`

	rows, err := r.Tx.QueryContext(r.Ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var scores []domain.BoardScoreDto
	for rows.Next() {
		var score domain.BoardScoreDto
		var nsNumber sql.NullInt64
		var ewNumber sql.NullInt64
		if err := rows.Scan(
			&score.DealNo,
			&score.NsTeamId, &nsNumber,
			&score.EwTeamId, &ewNumber,
			&score.NsPoints, &score.NsAward, &score.EwAward,
		); err != nil {
			return nil, err
		}
		if nsNumber.Valid {
			score.NsTeamNumber = int(nsNumber.Int64)
		}
		if ewNumber.Valid {
			score.EwTeamNumber = int(ewNumber.Int64)
		}
		scores = append(scores, score)
	}
	return scores, nil
}
