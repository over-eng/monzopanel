package eventstore

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/over-eng/monzopanel/libraries/models"
)

func (s *Store) InsertEvent(ctx context.Context, event *models.Event) error {
	properties, err := json.Marshal(event.Properties)
	if err != nil {
		return errors.Join(models.ErrInvalidEvent, err)
	}
	event.LoadedAt = time.Now()

	err = event.ValidateInsertable()
	if err != nil {
		return errors.Join(models.ErrInvalidEvent, err)
	}

	cql := `
		INSERT INTO events (
			id,
			event,
			team_id,
			distinct_id,
			properties,
			client_timestamp,
			created_at,
			loaded_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	return s.session.Query(
		cql,
		event.ID,
		event.Event,
		event.TeamID,
		event.DistinctID,
		properties,
		event.ClientTimestamp,
		event.CreatedAt,
		event.LoadedAt,
	).WithContext(ctx).Exec()
}
