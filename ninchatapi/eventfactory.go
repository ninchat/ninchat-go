package ninchatapi

import (
	"github.com/ninchat/ninchat-go"
)

// NewEvent creates an event object with the parameters specified by the
// clientEvent.  An UnexpectedEventError is returned if its type is unknown.
func NewEvent(clientEvent *ninchat.Event) (event Event, err error) {
	if clientEvent != nil {
		if f := EventFactories[clientEvent.String()]; f != nil {
			event = f()
			event.Init(clientEvent)
		} else {
			err = &UnexpectedEventError{clientEvent}
		}
	}
	return
}
