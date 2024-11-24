package eventstore

import (
	"encoding/base64"

	"github.com/over-eng/monzopanel/protos/event"
)

func generatePaginationToken(event *event.Event) string {
	return base64.StdEncoding.EncodeToString([]byte(event.CreatedAt.String()))
}

func DecodePaginationToken(token string) (string, error) {
	if token == "" {
		return "", nil
	}
	data, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", err
	}
	return string(data), nil
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
	token, err := DecodePaginationToken(paginationToken)
	if err != nil {
		return ListEventsByDistinctIDResult{}, err
	}

	var query string
	var args []any

	// query for one more than the request, to determine if more events exist.
	limit := pageSize + 1

	if token != "" {
		query = `
			SELECT
				id,
				event,
				team_id,
				distinct_id,
				properties,
				client_timestamp,
				created_at,
				loaded_at
			FROM events 
			WHERE team_id = ? 
				AND distinct_id = ? 
				AND created_at < ? 
			LIMIT ?
		`
		args = []any{teamID, distinctID, token, limit}
	} else {
		query = `
			SELECT
				id,
				event,
				team_id,
				distinct_id,
				properties,
				client_timestamp,
				created_at,
				loaded_at
			FROM events 
			WHERE team_id = ? 
				AND distinct_id = ? 
			LIMIT ?
		`
		args = []any{teamID, distinctID, limit}
	}

	iter := s.session.Query(query, args...).Iter()
	defer iter.Close()

	var events []*event.Event
	var hasMore bool

	// Scan results
	for i := 0; i < limit; i++ {
		event := &event.Event{}
		if !iter.Scan(
			&event.Id,
			&event.Event,
			&event.TeamId,
			&event.DistinctId,
			&event.Properties,
			&event.ClientTimestamp,
			&event.CreatedAt,
			&event.LoadedAt,
		) {
			break
		}
		if i < pageSize {
			events = append(events, event)
		} else {
			hasMore = true
			break
		}
	}

	// Generate next pagination token if there are more results
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
