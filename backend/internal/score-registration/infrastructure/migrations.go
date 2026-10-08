package score_registration

import (
	"database/sql"
)

func Migrate(db *sql.DB) {
	m0001_initial(db)
}

func m0001_initial(db *sql.DB) {
	_, err := db.Exec(`
		CREATE SCHEMA IF NOT EXISTS score_registration;

		CREATE TABLE IF NOT EXISTS score_registration.played_result (
			tournament_id UUID NOT NULL,
			deal_no INTEGER NOT NULL,
			ns_team_id UUID NOT NULL,
			ew_team_id UUID NOT NULL,
			contract TEXT NOT NULL,
			tricks INTEGER NOT NULL,
			declarer TEXT NOT NULL,
			ns_vulnerable BOOLEAN NOT NULL,
			ew_vulnerable BOOLEAN NOT NULL,

			PRIMARY KEY (tournament_id, deal_no, ns_team_id, ew_team_id)
		);

		CREATE TABLE IF NOT EXISTS score_registration.tournament_score (
			tournament_id UUID NOT NULL,
			method TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

			PRIMARY KEY (tournament_id, method)
		);

		CREATE TABLE IF NOT EXISTS score_registration.standing (
			tournament_id UUID NOT NULL,
			method TEXT NOT NULL,
			team_id UUID NOT NULL,
			total_score DOUBLE PRECISION NOT NULL,
			percentage DOUBLE PRECISION NOT NULL,
			rank INTEGER NOT NULL,

			PRIMARY KEY (tournament_id, method, team_id),
			FOREIGN KEY (tournament_id, method)
				REFERENCES score_registration.tournament_score (tournament_id, method)
				ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS score_registration.board_score (
			tournament_id UUID NOT NULL,
			method TEXT NOT NULL,
			deal_no INTEGER NOT NULL,
			ns_team_id UUID NOT NULL,
			ew_team_id UUID NOT NULL,
			ns_points INTEGER NOT NULL,
			ns_award DOUBLE PRECISION NOT NULL,
			ew_award DOUBLE PRECISION NOT NULL,

			PRIMARY KEY (tournament_id, method, deal_no, ns_team_id, ew_team_id),
			FOREIGN KEY (tournament_id, method)
				REFERENCES score_registration.tournament_score (tournament_id, method)
				ON DELETE CASCADE
		);
	`)
	if err != nil {
		panic(err)
	}
}
