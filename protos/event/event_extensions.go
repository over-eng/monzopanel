// this file wasn't generated, these are custom extensions

package event

import (
	"errors"
)

var ErrInvalidEvent = errors.New("event is not valid")

func (e *Event) ValidateQueueable() error {
	if e.Id == "" {
		return errors.Join(ErrInvalidEvent, errors.New("id"))
	}

	if e.Event == "" {
		return errors.Join(ErrInvalidEvent, errors.New("missing event name"))
	}

	if e.TeamId == "" {
		return errors.Join(ErrInvalidEvent, errors.New("missing team id"))
	}

	if e.ClientTimestamp == nil {
		return errors.Join(ErrInvalidEvent, errors.New("client timestamp not defined"))
	}
	return nil
}

// TODO: actually spec this
func (e *Event) ValidateInsertable() error {
	return e.ValidateQueueable()
}
