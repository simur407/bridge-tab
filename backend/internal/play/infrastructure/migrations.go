package play

import (
	"database/sql"
)

func Migrate(db *sql.DB) {
	_, err := db.Exec(`
		CREATE SCHEMA IF NOT EXISTS play;

		CREATE TABLE IF NOT EXISTS play.board_result (
			play_session_id UUID NOT NULL,
			deal_no INTEGER NOT NULL,
			ns_team_id UUID NOT NULL,
			ew_team_id UUID NOT NULL,
			contract TEXT NOT NULL,
			tricks INTEGER NOT NULL,
			declarer TEXT NOT NULL,
			opening_lead TEXT NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

			PRIMARY KEY (play_session_id, deal_no, ns_team_id, ew_team_id)
		);
	`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS play.team (
			play_session_id UUID NOT NULL,
			team_id UUID NOT NULL,
			number INTEGER NOT NULL,

			PRIMARY KEY (play_session_id, team_id)
		);

		CREATE TABLE IF NOT EXISTS play.team_player (
			play_session_id UUID NOT NULL,
			team_id UUID NOT NULL,
			player_id UUID NOT NULL,

			PRIMARY KEY (play_session_id, team_id, player_id)
		);
	`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS play.table_schedule (
			play_session_id UUID NOT NULL,
			round_no INTEGER NOT NULL,
			table_number INTEGER NOT NULL,
			ns_team_id UUID NOT NULL,
			ew_team_id UUID NOT NULL,
			set_id UUID NOT NULL,
			set_label TEXT NOT NULL,

			PRIMARY KEY (play_session_id, round_no, table_number)
		);

		CREATE TABLE IF NOT EXISTS play.table_schedule_board (
			play_session_id UUID NOT NULL,
			round_no INTEGER NOT NULL,
			table_number INTEGER NOT NULL,
			deal_no INTEGER NOT NULL,
			vulnerable TEXT NOT NULL,
			position INTEGER NOT NULL,

			PRIMARY KEY (play_session_id, round_no, table_number, deal_no)
		);

		DROP TABLE IF EXISTS play.match;
	`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		DO $$
		DECLARE
			play_table text;
		BEGIN
			FOREACH play_table IN ARRAY ARRAY[
				'board_result',
				'team',
				'team_player',
				'table_schedule',
				'table_schedule_board'
			]
			LOOP
				IF EXISTS (
					SELECT 1
					FROM information_schema.columns
					WHERE table_schema = 'play'
						AND information_schema.columns.table_name = play_table
						AND column_name = 'tournament_id'
				) THEN
					EXECUTE format(
						'ALTER TABLE play.%I RENAME COLUMN tournament_id TO play_session_id',
						play_table
					);
				END IF;
			END LOOP;
		END $$;
	`)
	if err != nil {
		panic(err)
	}
}
