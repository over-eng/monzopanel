package eventstore

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/over-eng/monzopanel/libraries/models"
	"github.com/over-eng/monzopanel/protos/event"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) InsertEvent(ctx context.Context, event *event.Event) error {
	properties, err := json.Marshal(event.Properties)
	if err != nil {
		return errors.Join(models.ErrInvalidEvent, err)
	}
	event.LoadedAt = timestamppb.Now()

	err = event.ValidateInsertable()
	if err != nil {
		return errors.Join(models.ErrInvalidEvent, err)
	}

	err = s.incrementEventCounterTable(ctx, event)
	if err != nil {
		s.log.Err(err).Msg("failed to increment event counter table")
		return errors.Join(errors.New("failed to increment event counter table"), err)
	}

	insertQuery := `
	INSERT INTO events_by_distinct_id (
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

	err = s.session.Query(
		insertQuery,
		event.Id,
		event.Event,
		event.TeamId,
		event.DistinctId,
		properties,
		event.ClientTimestamp.AsTime(),
		event.CreatedAt.AsTime(),
		event.LoadedAt.AsTime(),
	).WithContext(ctx).Exec()
	if err != nil {
		// on failure we need to undo the increase to the counter table
		decrementErr := s.decrementEventCounterTable(ctx, event)
		if decrementErr != nil {
			errStr := "increment successful, insert and decrement unsuccessful, data may be in an inconsistent state"
			s.log.Err(decrementErr).Msg(errStr)
			return errors.Join(errors.New(errStr), decrementErr)
		}

		return errors.Join(errors.New("failed to insert event"), err)
	}

	return nil
}

func (s *Store) incrementEventCounterTable(ctx context.Context, event *event.Event) error {
	bucketHour := event.CreatedAt.AsTime().Truncate(time.Hour)
	updateQuery := `
	UPDATE events_by_hour_counter
	SET event_count = event_count + 1
	WHERE
		team_id = ?
		AND distinct_id = ?
		AND bucket_hour = ?
		AND event = ?
	`
	return s.session.Query(
		updateQuery,
		event.TeamId,
		event.DistinctId,
		bucketHour,
		event.Event,
	).WithContext(ctx).Exec()
}

func (s *Store) decrementEventCounterTable(ctx context.Context, event *event.Event) error {
	bucketHour := event.CreatedAt.AsTime().Truncate(time.Hour)
	updateQuery := `
	UPDATE events_by_hour_counter
	SET event_count = event_count - 1
	WHERE
		team_id = ?
		AND distinct_id = ?
		AND bucket_hour = ?
		AND event = ?
	`
	return s.session.Query(
		updateQuery,
		event.TeamId,
		event.DistinctId,
		bucketHour,
		event.Event,
	).WithContext(ctx).Exec()
}
