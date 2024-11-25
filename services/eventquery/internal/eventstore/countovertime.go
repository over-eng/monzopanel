package eventstore

import (
	"context"
	"time"
)

const MAX_HOURS = 10000

type CountOvertime struct {
	Timestamp time.Time
	Count     int
}

func (s *Store) GetEventCountOvertime(
	ctx context.Context,
	teamID string,
	from time.Time,
	to time.Time,
) ([]*CountOvertime, error) {
	query := `
		SELECT
			event_count,
			bucket_hour
		FROM events_by_hour_counter
		WHERE
			team_id = ?
			AND bucket_hour >= ?
			AND bucket_hour <= ?
		LIMIT ?
		ALLOW FILTERING
	`

	iter := s.session.Query(query, teamID, from, to, MAX_HOURS).Iter()

	results := []*CountOvertime{}
	for {
		row := &CountOvertime{}
		if !iter.Scan(&row.Count, &row.Timestamp) {
			break
		}

		results = append(results, row)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return results, nil
}
