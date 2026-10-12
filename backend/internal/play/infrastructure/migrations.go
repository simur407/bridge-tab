package play

import (
	"database/sql"
)

func Migrate(db *sql.DB) {
	_, err := db.Exec(`
		CREATE SCHEMA IF NOT EXISTS play;

		CREATE TABLE IF NOT EXISTS play.board_result (
			tournament_id UUID NOT NULL,
			deal_no INTEGER NOT NULL,
			ns_team_id UUID NOT NULL,
			ew_team_id UUID NOT NULL,
			contract TEXT NOT NULL,
			tricks INTEGER NOT NULL,
			declarer TEXT NOT NULL,
			opening_lead TEXT NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

			PRIMARY KEY (tournament_id, deal_no, ns_team_id, ew_team_id)
		);
	`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS play.team (
			tournament_id UUID NOT NULL,
			team_id UUID NOT NULL,
			number INTEGER NOT NULL,

			PRIMARY KEY (tournament_id, team_id)
		);

		CREATE TABLE IF NOT EXISTS play.team_player (
			tournament_id UUID NOT NULL,
			team_id UUID NOT NULL,
			player_id UUID NOT NULL,

			PRIMARY KEY (tournament_id, team_id, player_id)
		);

		CREATE TABLE IF NOT EXISTS play.match (
			tournament_id UUID NOT NULL,
			deal_no INTEGER NOT NULL,
			ns_team_id UUID NOT NULL,
			ew_team_id UUID NOT NULL,
			vulnerable TEXT NOT NULL,

			PRIMARY KEY (tournament_id, deal_no, ns_team_id, ew_team_id)
		);
	`)
	if err != nil {
		panic(err)
	}
}
