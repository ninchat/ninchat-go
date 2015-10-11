package ninchatapi

import (
	"github.com/ninchat/ninchat-go"
)

// EventError wraps an "error" event or an unexpected event.
type EventError struct {
	Event *ninchat.Event
}

// String returns the event type.
func (err *EventError) String() string {
	return err.Event.String()
}

// Error returns a detailed error description.
func (err *EventError) Error() (s string) {
	switch name := err.Event.String(); name {
	case "error":
		s, _ = err.Event.Str("error_type")
		if r, ok := err.Event.Str("error_reason"); ok {
			s += ": " + r
		}

	default:
		s = "unexpected event type: " + name
	}

	return
}
