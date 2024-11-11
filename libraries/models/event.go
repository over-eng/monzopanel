package models

import (
	"errors"
	"time"
)

type Event struct {
	ID    string `json:"id"`
	Event string `json:"event"`

	TeamID     string `json:"team_id"`
	DistinctID string `json:"distinct_id"`

	// Event properties to me marshalled to a json string
	Properties map[string]interface{} `json:"properties"`

	// time as declared by the client
	ClientTimestamp time.Time `json:"client_timestamp"`
	// when the event first reached our servers
	CreatedAt time.Time `json:"created_at"`
	// when the event was inserted into our db
	LoadedAt time.Time `json:"loaded_at"`
}

var ErrInvalidEvent = errors.New("event is not valid")

func (e Event) ValidateQueueable() error {
	if e.ID == "" {
		return errors.Join(ErrInvalidEvent, errors.New("id"))
	}

	if e.Event == "" {
		return errors.Join(ErrInvalidEvent, errors.New("missing event name"))
	}

	if e.TeamID == "" {
		return errors.Join(ErrInvalidEvent, errors.New("missing team id"))
	}

	if e.ClientTimestamp.IsZero() {
		return errors.Join(ErrInvalidEvent, errors.New("client timestamp not defined"))
	}
	return nil
}

// TODO: actually spec this
func (e Event) ValidateInsertable() error {
	return e.ValidateQueueable()
}
