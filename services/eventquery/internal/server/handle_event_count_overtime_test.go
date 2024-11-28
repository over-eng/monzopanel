package server_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/over-eng/monzopanel/protos/event"
	"github.com/over-eng/monzopanel/services/eventquery/internal/eventstore"
	"github.com/over-eng/monzopanel/services/eventquery/internal/server"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGenerateBucketArray(t *testing.T) {

	tests := []struct {
		name     string
		counts   []*eventstore.CountOvertime
		from     time.Time
		to       time.Time
		period   event.TimePeriod
		expected []*event.TimeBucket
	}{
		{
			name: "should pad with zeros",
			counts: []*eventstore.CountOvertime{
				{Count: 2, Timestamp: time.Date(2024, 11, 28, 15, 0, 0, 0, time.UTC)},
			},
			from:   time.Date(2024, 11, 28, 14, 0, 0, 0, time.UTC),
			to:     time.Date(2024, 11, 28, 18, 0, 0, 0, time.UTC),
			period: event.TimePeriod_HOUR,
			expected: []*event.TimeBucket{
				{Count: 0, Timestamp: timestamppb.New(time.Date(2024, 11, 28, 14, 0, 0, 0, time.UTC))},
				{Count: 2, Timestamp: timestamppb.New(time.Date(2024, 11, 28, 15, 0, 0, 0, time.UTC))},
				{Count: 0, Timestamp: timestamppb.New(time.Date(2024, 11, 28, 16, 0, 0, 0, time.UTC))},
				{Count: 0, Timestamp: timestamppb.New(time.Date(2024, 11, 28, 17, 0, 0, 0, time.UTC))},
				{Count: 0, Timestamp: timestamppb.New(time.Date(2024, 11, 28, 18, 0, 0, 0, time.UTC))},
			},
		},
		{
			name: "should aggregate time periods",
			counts: []*eventstore.CountOvertime{
				{Count: 1, Timestamp: time.Date(2024, 11, 27, 14, 0, 0, 0, time.UTC)},
				{Count: 2, Timestamp: time.Date(2024, 11, 28, 15, 0, 0, 0, time.UTC)},
				{Count: 5, Timestamp: time.Date(2024, 11, 28, 16, 0, 0, 0, time.UTC)},
				{Count: 0, Timestamp: time.Date(2024, 11, 29, 17, 0, 0, 0, time.UTC)},
				{Count: 4, Timestamp: time.Date(2024, 11, 30, 18, 0, 0, 0, time.UTC)},
			},
			from:   time.Date(2024, 11, 27, 0, 0, 0, 0, time.UTC),
			to:     time.Date(2024, 11, 30, 0, 0, 0, 0, time.UTC),
			period: event.TimePeriod_DAY,
			expected: []*event.TimeBucket{
				{Count: 1, Timestamp: timestamppb.New(time.Date(2024, 11, 27, 0, 0, 0, 0, time.UTC))},
				{Count: 7, Timestamp: timestamppb.New(time.Date(2024, 11, 28, 0, 0, 0, 0, time.UTC))},
				{Count: 0, Timestamp: timestamppb.New(time.Date(2024, 11, 29, 0, 0, 0, 0, time.UTC))},
				{Count: 4, Timestamp: timestamppb.New(time.Date(2024, 11, 30, 0, 0, 0, 0, time.UTC))},
			},
		},
		{
			name: "inputs should also be truncated to nearest bucket",
			counts: []*eventstore.CountOvertime{
				{Count: 1, Timestamp: time.Date(2024, 11, 27, 14, 0, 0, 0, time.UTC)},
				{Count: 2, Timestamp: time.Date(2024, 11, 28, 15, 0, 0, 0, time.UTC)},
				{Count: 5, Timestamp: time.Date(2024, 11, 28, 16, 0, 0, 0, time.UTC)},
				{Count: 0, Timestamp: time.Date(2024, 11, 29, 17, 0, 0, 0, time.UTC)},
				{Count: 4, Timestamp: time.Date(2024, 11, 30, 18, 0, 0, 0, time.UTC)},
			},
			from:   time.Date(2024, 11, 27, 14, 0, 0, 0, time.UTC),
			to:     time.Date(2024, 11, 30, 8, 0, 0, 0, time.UTC),
			period: event.TimePeriod_DAY,
			expected: []*event.TimeBucket{
				{Count: 1, Timestamp: timestamppb.New(time.Date(2024, 11, 27, 0, 0, 0, 0, time.UTC))},
				{Count: 7, Timestamp: timestamppb.New(time.Date(2024, 11, 28, 0, 0, 0, 0, time.UTC))},
				{Count: 0, Timestamp: timestamppb.New(time.Date(2024, 11, 29, 0, 0, 0, 0, time.UTC))},
				{Count: 4, Timestamp: timestamppb.New(time.Date(2024, 11, 30, 0, 0, 0, 0, time.UTC))},
			},
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actual, err := server.GenerateBucketArray(tc.counts, tc.from, tc.to, tc.period)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tc.expected, actual) {
				t.Fatalf("expected: %v, got: %v", tc.expected, actual)
			}
		})
	}
}
