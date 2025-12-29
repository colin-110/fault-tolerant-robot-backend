package storage

import "database/sql"

func InitSchema(db *sql.DB) error {
	statements := []string{
		`
		CREATE TABLE IF NOT EXISTS commands (
			command_id TEXT PRIMARY KEY,
			robot_id TEXT NOT NULL,
			state TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS robot_state (
			robot_id TEXT PRIMARY KEY,
			last_executed_command_id TEXT
		);
		`,
	}

	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}
	return nil
}
