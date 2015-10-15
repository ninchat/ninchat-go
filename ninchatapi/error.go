package ninchatapi

import (
	"github.com/ninchat/ninchat-go"
)

// Error describes an "error" event.
func (event *Error) Error() (s string) {
	s = event.ErrorType
	if event.ErrorReason != nil {
		s += ": " + *event.ErrorReason
	}
	return
}

// newRequestMalformedError synthesizes an "error" event.
func newRequestMalformedError(reason string) *Error {
	return &Error{
		ErrorType:   "request_malformed",
		ErrorReason: &reason,
	}
}

// UnexpectedEventError wraps an event and implements the error interface.
type UnexpectedEventError struct {
	Event *ninchat.Event
}

// Error returns an error description.
func (e *UnexpectedEventError) Error() string {
	return "unexpected event type: " + e.Event.String()
}

// String returns an error description.
func (e *UnexpectedEventError) String() string {
	return e.Error()
}
