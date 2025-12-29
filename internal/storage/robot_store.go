package storage

import "database/sql"

type RobotStore struct {
	db *sql.DB
}

func NewRobotStore(db *sql.DB) *RobotStore {
	return &RobotStore{db: db}
}

func (s *RobotStore) GetLastExecuted(robotID string) (string, error) {
	row := s.db.QueryRow(
		`SELECT last_executed_command_id
		 FROM robot_state WHERE robot_id = ?`,
		robotID,
	)

	var cmdID string
	if err := row.Scan(&cmdID); err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return cmdID, nil
}

func (s *RobotStore) UpdateLastExecuted(robotID, commandID string) error {
	_, err := s.db.Exec(
		`
		INSERT INTO robot_state (robot_id, last_executed_command_id)
		VALUES (?, ?)
		ON CONFLICT(robot_id)
		DO UPDATE SET last_executed_command_id = excluded.last_executed_command_id
		`,
		robotID, commandID,
	)
	return err
}
