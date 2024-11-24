package eventstore

import (
	"encoding/base64"
	"time"

	"github.com/over-eng/monzopanel/protos/event"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func generatePaginationToken(event *event.Event) string {
	createdAt := event.CreatedAt.AsTime()
	return base64.StdEncoding.EncodeToString([]byte(createdAt.Format(time.RFC1123)))
}

func decodePaginationToken(token string) (time.Time, error) {
	if token == "" {
		return time.Time{}, nil
	}
	data, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return time.Time{}, err
	}

	decoded, err := time.Parse(time.RFC1123, string(data))
	if err != nil {
		return time.Time{}, err
	}

	return decoded, nil
}

type ListEventsByDistinctIDResult struct {
	Events              []*event.Event
	NextPaginationToken string
}

func (s *Store) ListEventsByDistinctID(
	teamID string,
	distinctID string,
	pageSize int,
	paginationToken string,
) (ListEventsByDistinctIDResult, error) {
	createdAtToken, err := decodePaginationToken(paginationToken)
	if err != nil {
		return ListEventsByDistinctIDResult{}, err
	}

	var query string
	var args []any

	// query for one more than the request, to determine if more events exist.
	limit := pageSize + 1

	if !createdAtToken.IsZero() {
		query = `
			SELECT
				id,
				event,
				team_id,
				distinct_id,
				client_timestamp,
				created_at,
				loaded_at,
				properties
			FROM events_by_distinct_id
			WHERE
				team_id = ? 
				AND distinct_id = ? 
				AND created_at < ? 
			LIMIT ?
		`
		args = []any{teamID, distinctID, createdAtToken, limit}
	} else {
		query = `
			SELECT
				id,
				event,
				team_id,
				distinct_id,
				client_timestamp,
				created_at,
				loaded_at,
				properties
			FROM events_by_distinct_id 
			WHERE
				team_id = ? 
				AND distinct_id = ? 
			LIMIT ?
		`
		args = []any{teamID, distinctID, limit}
	}

	iter := s.session.Query(query, args...).Iter()
	defer iter.Close()

	var events []*event.Event
	var hasMore bool
	for i := 0; i < limit; i++ {
		event := &event.Event{}

		// it's not possible to natively scan the protobuf struct
		// on these fields so we process them separately
		var (
			clientTimestamp time.Time
			createdAt       time.Time
			loadedAt        time.Time
			properties      string
		)

		if !iter.Scan(
			&event.Id,
			&event.Event,
			&event.TeamId,
			&event.DistinctId,
			&clientTimestamp,
			&createdAt,
			&loadedAt,
			&properties,
		) {
			break
		}
		if i < pageSize {

			event.ClientTimestamp = timestamppb.New(clientTimestamp)
			event.CreatedAt = timestamppb.New(createdAt)
			event.LoadedAt = timestamppb.New(loadedAt)

			// initialise event.Properties to prevent nil dereference
			event.Properties = &structpb.Struct{}
			if properties != "" {
				err := protojson.Unmarshal([]byte(properties), event.Properties)
				if err != nil {
					s.log.Err(err).Msg("Failed to unmarshal properties")
					return ListEventsByDistinctIDResult{}, err
				}
			}

			events = append(events, event)
		} else {
			hasMore = true
			break
		}
	}

	var nextToken string
	if hasMore && len(events) > 0 {
		nextToken = generatePaginationToken(events[len(events)-1])
	}

	result := ListEventsByDistinctIDResult{
		Events:              events,
		NextPaginationToken: nextToken,
	}

	return result, nil
}
