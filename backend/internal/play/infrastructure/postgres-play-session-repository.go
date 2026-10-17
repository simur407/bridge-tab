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

func (r *PostgresPlaySessionRepository) Load(playSessionId domain.PlaySessionId) (*domain.PlaySession, error) {
	var session domain.PlaySession

	teamRows, err := r.Tx.QueryContext(r.Ctx, `
		SELECT team_id, number
		FROM play.team
		WHERE play_session_id = $1`, playSessionId)
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
		WHERE play_session_id = $1`, playSessionId)
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
		SELECT round_no, table_number, ns_team_id, ew_team_id, set_id, set_label
		FROM play.table_schedule
		WHERE play_session_id = $1
		ORDER BY round_no, table_number`, playSessionId)
	if err != nil {
		return nil, err
	}
	defer matchRows.Close()

	for matchRows.Next() {
		var schedule domain.TableSchedule
		if err := matchRows.Scan(&schedule.RoundNo, &schedule.TableNumber, &schedule.NsTeamId, &schedule.EwTeamId, &schedule.SetId, &schedule.SetLabel); err != nil {
			return nil, err
		}
		session.State.TableSchedules = append(session.State.TableSchedules, &schedule)
	}
	if err := matchRows.Err(); err != nil {
		return nil, err
	}

	boardRows, err := r.Tx.QueryContext(r.Ctx, `
		SELECT round_no, table_number, deal_no, vulnerable
		FROM play.table_schedule_board
		WHERE play_session_id = $1
		ORDER BY round_no, table_number, position`, playSessionId)
	if err != nil {
		return nil, err
	}
	defer boardRows.Close()

	for boardRows.Next() {
		var roundNo int
		var tableNumber int
		var board domain.TableScheduleBoard
		if err := boardRows.Scan(&roundNo, &tableNumber, &board.DealNo, &board.Vulnerable); err != nil {
			return nil, err
		}
		index := slices.IndexFunc(session.State.TableSchedules, func(schedule *domain.TableSchedule) bool {
			return schedule.RoundNo == roundNo && schedule.TableNumber == tableNumber
		})
		if index == -1 {
			return nil, errors.New("play is not configured correctly")
		}
		session.State.TableSchedules[index].Boards = append(session.State.TableSchedules[index].Boards, board)
	}
	if err := boardRows.Err(); err != nil {
		return nil, err
	}

	if len(session.State.Teams) == 0 && len(session.State.TableSchedules) == 0 {
		return nil, domain.ErrPlaySessionNotProvisioned
	}

	session.State.PlaySessionId = playSessionId
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
			INSERT INTO play.team (play_session_id, team_id, number)
			VALUES ($1, $2, $3)`,
			event.PlaySessionId, team.Id, team.Number,
		)
		if err != nil {
			return err
		}
		for _, player := range team.Players {
			_, err = r.Tx.ExecContext(r.Ctx, `
				INSERT INTO play.team_player (play_session_id, team_id, player_id)
				VALUES ($1, $2, $3)`,
				event.PlaySessionId, team.Id, player,
			)
			if err != nil {
				return err
			}
		}
	}

	for _, schedule := range event.TableSchedules {
		_, err := r.Tx.ExecContext(r.Ctx, `
			INSERT INTO play.table_schedule (play_session_id, round_no, table_number, ns_team_id, ew_team_id, set_id, set_label)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			event.PlaySessionId, schedule.RoundNo, schedule.TableNumber, schedule.NsTeamId, schedule.EwTeamId, schedule.SetId, schedule.SetLabel,
		)
		if err != nil {
			return err
		}
		for position, board := range schedule.Boards {
			_, err = r.Tx.ExecContext(r.Ctx, `
				INSERT INTO play.table_schedule_board (play_session_id, round_no, table_number, deal_no, vulnerable, position)
				VALUES ($1, $2, $3, $4, $5, $6)`,
				event.PlaySessionId, schedule.RoundNo, schedule.TableNumber, board.DealNo, board.Vulnerable, position,
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
