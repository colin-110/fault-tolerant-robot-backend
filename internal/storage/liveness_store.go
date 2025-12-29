package storage

import (
	"database/sql"
	"time"
)

type LivenessStore struct {
	db *sql.DB
}

func NewLivenessStore(db *sql.DB) *LivenessStore {
	return &LivenessStore{db: db}
}

func (s *LivenessStore) Heartbeat(robotID string) error {
	_, err := s.db.Exec(`
		INSERT INTO robot_liveness (robot_id, last_seen_unix_ms)
		VALUES (?, ?)
		ON CONFLICT(robot_id)
		DO UPDATE SET last_seen_unix_ms = excluded.last_seen_unix_ms
	`,
		robotID,
		time.Now().UnixMilli(),
	)
	return err
}

func (s *LivenessStore) Expired(threshold time.Duration) ([]string, error) {
	cutoff := time.Now().Add(-threshold).UnixMilli()

	rows, err := s.db.Query(`
		SELECT robot_id FROM robot_liveness
		WHERE last_seen_unix_ms < ?
	`, cutoff)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dead []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		dead = append(dead, id)
	}
	return dead, nil
}
