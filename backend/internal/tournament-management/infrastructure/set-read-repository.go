package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
	"context"
	"database/sql"
)

type PostgresSetReadRepository struct {
	Ctx context.Context
	Tx  *sql.Tx
}

func (r *PostgresSetReadRepository) FindAll(tournamentId *string) ([]domain.SetDto, error) {
	setRows, err := r.Tx.QueryContext(r.Ctx, `
		SELECT id, label
		FROM tournament_management.tournament_set
		WHERE tournament_id = $1
		ORDER BY label`, tournamentId)
	if err != nil {
		return nil, err
	}

	var sets []domain.SetDto
	setIndex := map[string]int{}
	for setRows.Next() {
		var set domain.SetDto
		if err := setRows.Scan(&set.Id, &set.Label); err != nil {
			return nil, err
		}
		set.BoardNos = []int{}
		setIndex[set.Id] = len(sets)
		sets = append(sets, set)
	}

	boardRows, err := r.Tx.QueryContext(r.Ctx, `
		SELECT set_id, board_no
		FROM tournament_management.board_protocol
		WHERE tournament_id = $1 AND set_id IS NOT NULL
		ORDER BY board_no`, tournamentId)
	if err != nil {
		return nil, err
	}

	for boardRows.Next() {
		var setId string
		var boardNo int
		if err := boardRows.Scan(&setId, &boardNo); err != nil {
			return nil, err
		}
		if i, ok := setIndex[setId]; ok {
			sets[i].BoardNos = append(sets[i].BoardNos, boardNo)
		}
	}

	return sets, nil
}
