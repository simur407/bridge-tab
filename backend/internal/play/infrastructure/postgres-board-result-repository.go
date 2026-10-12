package play

import (
	"context"
	"database/sql"
	"errors"

	domain "bridge-tab/internal/play/domain"
)

type PostgresBoardResultRepository struct {
	Ctx context.Context
	Tx  *sql.Tx
}

func (r *PostgresBoardResultRepository) Load(tournamentId domain.TournamentId, dealNo int, teamA domain.TeamId, teamB domain.TeamId) (*domain.BoardResult, error) {
	row := r.Tx.QueryRowContext(r.Ctx, `
		SELECT deal_no, ns_team_id, ew_team_id, contract, tricks, declarer, opening_lead
		FROM play.board_result
		WHERE tournament_id = $1 AND deal_no = $2
			AND ((ns_team_id = $3 AND ew_team_id = $4) OR (ns_team_id = $4 AND ew_team_id = $3))`,
		tournamentId, dealNo, teamA, teamB)

	var deal int
	var nsTeamId, ewTeamId domain.TeamId
	var contract, declarer, openingLead string
	var tricks int
	if err := row.Scan(&deal, &nsTeamId, &ewTeamId, &contract, &tricks, &declarer, &openingLead); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrBoardResultNotFound
		}
		return nil, err
	}

	return &domain.BoardResult{
		State: domain.BoardResultState{
			TournamentId: tournamentId,
			DealNo:       deal,
			NsTeamId:     nsTeamId,
			EwTeamId:     ewTeamId,
			Contract:     contract,
			Tricks:       tricks,
			Declarer:     declarer,
			OpeningLead:  openingLead,
		},
	}, nil
}

func (r *PostgresBoardResultRepository) Save(result *domain.BoardResult) error {
	for _, event := range result.GetEvents() {
		var err error
		switch event := event.(type) {
		case domain.BoardResultRecorded:
			err = r.boardResultRecorded(event)
		case domain.BoardResultEdited:
			err = r.boardResultEdited(event)
		default:
			return errors.New("unknown event")
		}
		if err != nil {
			return err
		}
	}
	result.Commit()
	return nil
}

func (r *PostgresBoardResultRepository) boardResultRecorded(event domain.BoardResultRecorded) error {
	_, err := r.Tx.ExecContext(r.Ctx, `
		INSERT INTO play.board_result (
			tournament_id, deal_no, ns_team_id, ew_team_id,
			contract, tricks, declarer, opening_lead
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		event.TournamentId, event.DealNo, event.NsTeamId, event.EwTeamId,
		event.Contract, event.Tricks, event.Declarer, event.OpeningLead,
	)
	return err
}

func (r *PostgresBoardResultRepository) boardResultEdited(event domain.BoardResultEdited) error {
	_, err := r.Tx.ExecContext(r.Ctx, `
		UPDATE play.board_result
		SET contract = $5, tricks = $6, declarer = $7, opening_lead = $8, updated_at = now()
		WHERE tournament_id = $1 AND deal_no = $2 AND ns_team_id = $3 AND ew_team_id = $4`,
		event.TournamentId, event.DealNo, event.NsTeamId, event.EwTeamId,
		event.Contract, event.Tricks, event.Declarer, event.OpeningLead,
	)
	return err
}
