package storage

import (
	"database/sql"

	"github.com/colin-110/fault-tolerant-robot-backend/internal/domain"
)

type CommandStore struct {
	db *sql.DB
}

func NewCommandStore(db *sql.DB) *CommandStore {
	return &CommandStore{db: db}
}

func (s *CommandStore) Create(cmd domain.Command) error {
	_, err := s.db.Exec(
		`INSERT INTO commands (command_id, robot_id, state)
		 VALUES (?, ?, ?)`,
		cmd.CommandID, cmd.RobotID, cmd.State,
	)
	return err
}

func (s *CommandStore) Get(commandID string) (*domain.Command, error) {
	row := s.db.QueryRow(
		`SELECT command_id, robot_id, state
		 FROM commands WHERE command_id = ?`,
		commandID,
	)

	var cmd domain.Command
	if err := row.Scan(&cmd.CommandID, &cmd.RobotID, &cmd.State); err != nil {
		return nil, err
	}
	return &cmd, nil
}

func (s *CommandStore) UpdateState(commandID string, next domain.CommandState) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	row := tx.QueryRow(
		`SELECT state FROM commands WHERE command_id = ?`,
		commandID,
	)

	var current domain.CommandState
	if err := row.Scan(&current); err != nil {
		return err
	}

	cmd := domain.Command{State: current}
	if err := domain.Transition(&cmd, next); err != nil {
		return err
	}

	if _, err := tx.Exec(
		`UPDATE commands
		 SET state = ?, updated_at = CURRENT_TIMESTAMP
		 WHERE command_id = ?`,
		next, commandID,
	); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *CommandStore) FailInFlightCommands(robotID string) error {
	_, err := s.db.Exec(`
		UPDATE commands
		SET state = 'FAILED',
		    updated_at = CURRENT_TIMESTAMP
		WHERE robot_id = ?
		  AND state IN ('SENT', 'ACKED')
	`, robotID)
	return err
}
func (s *CommandStore) FailAllUnfinished() error {
	_, err := s.db.Exec(`
		UPDATE commands
		SET state = 'FAILED',
		    updated_at = CURRENT_TIMESTAMP
		WHERE state IN ('CREATED', 'SENT', 'ACKED')
	`)
	return err
}
