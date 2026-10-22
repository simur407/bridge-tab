package umpirecall

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const (
	SourceTablet     = "tablet"
	SourceContestant = "contestant"
	ComingWindow     = 30 * time.Second
)

type Call struct {
	Id           string
	TournamentId string
	TableNumber  int
	RoundNo      int
	NsTeamNumber int
	EwTeamNumber int
	Source       string
	CreatedAt    time.Time
	AnsweredAt   *time.Time
}

type Repository struct {
	Ctx context.Context
	Tx  *sql.Tx
}

const callColumns = `
	id, tournament_id, table_number, round_no,
	ns_team_number, ew_team_number, source, created_at, answered_at
`

// Insert stores a call when this tournament and table have no unanswered row.
// A duplicate open call returns inserted=false and no error.
func (r *Repository) Insert(call Call) (bool, error) {
	if call.Id == "" {
		call.Id = uuid.NewString()
	}
	var existing string
	err := r.Tx.QueryRowContext(r.Ctx, `
		SELECT id FROM umpire_call.call
		WHERE tournament_id = $1 AND table_number = $2 AND answered_at IS NULL
	`, call.TournamentId, call.TableNumber).Scan(&existing)
	if err == nil {
		return false, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}

	if _, err := r.Tx.ExecContext(r.Ctx, `SAVEPOINT umpire_call_insert`); err != nil {
		return false, err
	}
	_, err = r.Tx.ExecContext(r.Ctx, `
		INSERT INTO umpire_call.call (
			id, tournament_id, table_number, round_no,
			ns_team_number, ew_team_number, source
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, call.Id, call.TournamentId, call.TableNumber, call.RoundNo, call.NsTeamNumber, call.EwTeamNumber, call.Source)
	if isUniqueViolation(err) {
		if _, rbErr := r.Tx.ExecContext(r.Ctx, `ROLLBACK TO SAVEPOINT umpire_call_insert`); rbErr != nil {
			return false, rbErr
		}
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if _, err := r.Tx.ExecContext(r.Ctx, `RELEASE SAVEPOINT umpire_call_insert`); err != nil {
		return false, err
	}
	return true, nil
}

func (r *Repository) ListOpen(tournamentId string) ([]Call, error) {
	return r.query(`
		SELECT `+callColumns+`
		FROM umpire_call.call
		WHERE tournament_id = $1 AND answered_at IS NULL
		ORDER BY created_at DESC
	`, tournamentId)
}

func (r *Repository) List(tournamentId string) ([]Call, error) {
	return r.query(`
		SELECT `+callColumns+`
		FROM umpire_call.call
		WHERE tournament_id = $1
		ORDER BY (answered_at IS NULL) DESC, created_at DESC
	`, tournamentId)
}

func (r *Repository) LatestForTable(tournamentId string, tableNumber int) (*Call, error) {
	row := r.Tx.QueryRowContext(r.Ctx, `
		SELECT `+callColumns+`
		FROM umpire_call.call
		WHERE tournament_id = $1 AND table_number = $2
		ORDER BY created_at DESC
		LIMIT 1
	`, tournamentId, tableNumber)
	call, err := scanCall(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &call, nil
}

func (r *Repository) MarkAnswered(tournamentId, id string) (bool, error) {
	result, err := r.Tx.ExecContext(r.Ctx, `
		UPDATE umpire_call.call
		SET answered_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND tournament_id = $2 AND answered_at IS NULL
	`, id, tournamentId)
	if err != nil {
		return false, err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (r *Repository) query(query string, tournamentId string) ([]Call, error) {
	rows, err := r.Tx.QueryContext(r.Ctx, query, tournamentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	calls := make([]Call, 0)
	for rows.Next() {
		call, err := scanCall(rows)
		if err != nil {
			return nil, err
		}
		calls = append(calls, call)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return calls, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanCall(row scanner) (Call, error) {
	var call Call
	var answered sql.NullTime
	err := row.Scan(
		&call.Id,
		&call.TournamentId,
		&call.TableNumber,
		&call.RoundNo,
		&call.NsTeamNumber,
		&call.EwTeamNumber,
		&call.Source,
		&call.CreatedAt,
		&answered,
	)
	if err != nil {
		return Call{}, err
	}
	if answered.Valid {
		t := answered.Time
		call.AnsweredAt = &t
	}
	return call, nil
}

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23505"
}
