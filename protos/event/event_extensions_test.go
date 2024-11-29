package event_test

import (
	"testing"

	"github.com/over-eng/monzopanel/protos/event"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestValidateQueueable(t *testing.T) {

	tabletests := []struct {
		name  string
		event *event.Event
		error error
	}{
		{
			name: "valid event",
			event: &event.Event{
				Id:              "test-id",
				Event:           "test-event",
				TeamId:          "test-team-id",
				DistinctId:      "test-distinct-id",
				Properties:      &structpb.Struct{},
				ClientTimestamp: timestamppb.Now(),
			},
			error: nil,
		},
		{
			name: "valid minimum event",
			event: &event.Event{
				Id:              "test-id",
				Event:           "test-event",
				TeamId:          "test-team-id",
				DistinctId:      "test-distinct-id",
				ClientTimestamp: timestamppb.Now(),
			},
			error: nil,
		},
		{
			name: "missing id",
			event: &event.Event{
				Event:           "test-event",
				TeamId:          "test-team-id",
				DistinctId:      "test-distinct-id",
				Properties:      &structpb.Struct{},
				ClientTimestamp: timestamppb.Now(),
			},
			error: event.ErrInvalidEvent,
		},
		{
			name: "missing event name",
			event: &event.Event{
				Id:              "test-id",
				TeamId:          "test-team-id",
				DistinctId:      "test-distinct-id",
				Properties:      &structpb.Struct{},
				ClientTimestamp: timestamppb.Now(),
			},
			error: event.ErrInvalidEvent,
		},
		{
			name: "missing team id",
			event: &event.Event{
				Id:              "test-id",
				Event:           "test-event",
				DistinctId:      "test-distinct-id",
				Properties:      &structpb.Struct{},
				ClientTimestamp: timestamppb.Now(),
			},
			error: event.ErrInvalidEvent,
		},
		{
			name: "missing client timestamp",
			event: &event.Event{
				Id:         "test-id",
				Event:      "test-event",
				TeamId:     "test-team-id",
				DistinctId: "test-distinct-id",
				Properties: &structpb.Struct{},
			},
			error: event.ErrInvalidEvent,
		},
	}

	for _, test := range tabletests {
		t.Run(test.name, func(t *testing.T) {
			err := test.event.ValidateQueueable()
			assert.ErrorIs(t, err, test.error)
		})
	}
}
