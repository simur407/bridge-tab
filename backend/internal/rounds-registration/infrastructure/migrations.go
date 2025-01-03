package rounds_registration

import (
	"database/sql"
)

// constraint naming: https://stackoverflow.com/a/4108266

func Migrate(db *sql.DB) {
	m0001_initial(db)
	m0002_add_updated_at(db)
}

func m0001_initial(db *sql.DB) {
	_, err := db.Exec(`
		CREATE SCHEMA IF NOT EXISTS rounds_registration;

		CREATE TABLE IF NOT EXISTS rounds_registration.round (
			game_session_id UUID NOT NULL,
			deal_no INTEGER NOT NULL,
			ns_team_id UUID NOT NULL,
			ew_team_id UUID NOT NULL,
			contract TEXT,
			declarer TEXT,
			tricks INTEGER,
			opening_lead TEXT,
			updated_at TIMESTAMPTZ,

			CONSTRAINT round_game_session_deal_no_unique UNIQUE (game_session_id, deal_no, ns_team_id, ew_team_id)
		);

		CREATE TABLE IF NOT EXISTS rounds_registration.team_players (
			game_session_id UUID NOT NULL,
			team_id UUID NOT NULL,
			player_id UUID NOT NULL,

			CONSTRAINT team_players_game_session_team_player_unique UNIQUE (game_session_id, team_id, player_id)
		)
	`)
	if err != nil {
		panic(err)
	}
}

func m0002_add_updated_at(db *sql.DB) {
	_, err := db.Exec(`
		ALTER TABLE rounds_registration.round
		ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ;
	`)
	if err != nil {
		panic(err)
	}
}
