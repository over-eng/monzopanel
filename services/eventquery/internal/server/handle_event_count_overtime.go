package server

import (
	"context"
	"errors"
	"time"

	"github.com/over-eng/monzopanel/protos/event"
	"github.com/over-eng/monzopanel/services/eventquery/internal/eventstore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) EventCountOvertime(
	ctx context.Context,
	req *event.EventCountOvertimeRequest,
) (*event.EventCountOvertimeResponse, error) {
	if req.TeamId == "" {
		return nil, status.Error(codes.InvalidArgument, "team_id is required")
	}

	from := req.From.AsTime()
	to := req.To.AsTime()

	if from.After(to) {
		return nil, errors.New("from timestamp is before to timestamp")
	}

	counts, err := s.eventstore.GetEventCountOvertime(ctx, req.TeamId, from, to)
	if err != nil {
		s.log.Err(err).Msg("failed to get event counts overtime")
		return nil, err
	}

	buckets, err := GenerateBucketArray(counts, from, to, req.Period)
	if err != nil {
		s.log.Err(err).Msg("failed to summarise counts")
		return nil, err
	}
	response := &event.EventCountOvertimeResponse{
		Buckets: buckets,
	}
	return response, nil
}

// The raw buckets returned from the database can contain gaps,
// so they are filled in with a zero value.
// hour windows are also aggregated for larger time periods.
func GenerateBucketArray(
	counts []*eventstore.CountOvertime,
	from time.Time,
	to time.Time,
	period event.TimePeriod,
) ([]*event.TimeBucket, error) {

	// Determine the increment based on the period
	var increment time.Duration
	switch period {
	case event.TimePeriod_HOUR:
		increment = time.Hour
	case event.TimePeriod_DAY:
		increment = 24 * time.Hour
	case event.TimePeriod_WEEK:
		increment = 7 * 24 * time.Hour
	default:
		return nil, errors.New("invalid time period supplied")
	}

	// truncate inputs
	to = to.Truncate(increment)
	from = from.Truncate(increment)

	current := from
	countsIdx := 0
	countsLength := len(counts)
	buckets := []*event.TimeBucket{}
	for current.Before(to) || current.Equal(to) {

		var count int64
		for {
			if countsIdx > countsLength-1 {
				break
			}

			c := counts[countsIdx]
			t := c.Timestamp.Truncate(increment)
			if t != current {
				break
			}
			count += int64(c.Count)
			countsIdx++
		}

		bucket := &event.TimeBucket{
			Count:     count,
			Timestamp: timestamppb.New(current),
		}
		buckets = append(buckets, bucket)
		current = current.Add(increment)
	}

	return buckets, nil
}
