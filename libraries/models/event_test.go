package models_test

import (
	"testing"
	"time"

	"github.com/over-eng/monzopanel/libraries/models"
	"github.com/stretchr/testify/assert"
)

func TestValidateQueueable(t *testing.T) {

	tabletests := []struct {
		name  string
		event models.Event
		error error
	}{
		{
			name: "valid event",
			event: models.Event{
				ID:              "test-id",
				Event:           "test-event",
				TeamID:          "test-team-id",
				DistinctID:      "test-distinct-id",
				Properties:      make(map[string]interface{}),
				ClientTimestamp: time.Now(),
			},
			error: nil,
		},
		{
			name: "valid minimum event",
			event: models.Event{
				ID:              "test-id",
				Event:           "test-event",
				TeamID:          "test-team-id",
				ClientTimestamp: time.Now(),
			},
			error: nil,
		},
		{
			name: "missing id",
			event: models.Event{
				Event:           "test-event",
				TeamID:          "test-team-id",
				DistinctID:      "test-distinct-id",
				Properties:      make(map[string]interface{}),
				ClientTimestamp: time.Now(),
			},
			error: models.ErrInvalidEvent,
		},
		{
			name: "missing event name",
			event: models.Event{
				ID:              "test-id",
				TeamID:          "test-team-id",
				DistinctID:      "test-distinct-id",
				Properties:      make(map[string]interface{}),
				ClientTimestamp: time.Now(),
			},
			error: models.ErrInvalidEvent,
		},
		{
			name: "missing team id",
			event: models.Event{
				ID:              "test-id",
				Event:           "test-event",
				DistinctID:      "test-distinct-id",
				Properties:      make(map[string]interface{}),
				ClientTimestamp: time.Now(),
			},
			error: models.ErrInvalidEvent,
		},
		{
			name: "missing client timestamp",
			event: models.Event{
				ID:         "test-id",
				Event:      "test-event",
				TeamID:     "test-team-id",
				DistinctID: "test-distinct-id",
				Properties: make(map[string]interface{}),
			},
			error: models.ErrInvalidEvent,
		},
	}

	for _, test := range tabletests {
		t.Run(test.name, func(t *testing.T) {
			err := test.event.ValidateQueueable()
			assert.ErrorIs(t, err, test.error)
		})
	}
}
