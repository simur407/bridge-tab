package umpirecall

import "database/sql"

func Migrate(db *sql.DB) {
	_, err := db.Exec(`
		CREATE SCHEMA IF NOT EXISTS umpire_call;

		CREATE TABLE IF NOT EXISTS umpire_call.call (
			id UUID PRIMARY KEY,
			tournament_id UUID NOT NULL,
			table_number INTEGER NOT NULL,
			round_no INTEGER NOT NULL,
			ns_team_number INTEGER NOT NULL,
			ew_team_number INTEGER NOT NULL,
			source TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
			answered_at TIMESTAMPTZ
		);

		CREATE UNIQUE INDEX IF NOT EXISTS umpire_call_one_open_per_table
			ON umpire_call.call (tournament_id, table_number)
			WHERE answered_at IS NULL;

		CREATE INDEX IF NOT EXISTS umpire_call_tournament_created
			ON umpire_call.call (tournament_id, created_at DESC);
	`)
	if err != nil {
		panic(err)
	}
}
