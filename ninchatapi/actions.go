package ninchatapi

import (
	"github.com/ninchat/ninchat-go"
)

// Action interface is implemented by all action structs.
type Action interface {
	// String returns the action type.
	String() string

	newClientAction() (*ninchat.Action, error)
}

// Sender can send actions.  It is implemented by ninchat.Session and Caller.
type Sender interface {
	// Send an action.
	Send(action *ninchat.Action) error
}

// Call an action.  If sender is nil, a trivial Caller will be used.
func Call(sender Sender, events chan<- *ninchat.Event, action Action) (err error) {
	clientAction, err := action.newClientAction()
	if err != nil {
		close(events)
		return
	}

	clientAction.OnReply = func(e *ninchat.Event) {
		if e == nil {
			close(events)
		} else {
			events <- e
			if e.LastReply {
				close(events)
			}
		}
	}

	if sender == nil {
		sender = &defaultCaller
	}

	err = sender.Send(clientAction)
	if err != nil {
		close(events)
	}
	return
}

// Send an action.
func Send(session *ninchat.Session, action Action) (err error) {
	clientAction, err := action.newClientAction()
	if err != nil {
		return
	}

	session.Send(clientAction)
	return
}

func unaryCall(sender Sender, action Action, event Event) (ok bool, err error) {
	c := make(chan *ninchat.Event, 1) // XXX: why doesn't this work without buffering?

	if err = Call(sender, c, action); err != nil {
		return
	}

	clientEvent := <-c
	flush(c)

	if clientEvent == nil {
		return
	}

	ok = true

	if clientEvent.String() == "error" {
		err = NewError(clientEvent)
	} else {
		err = event.Init(clientEvent)
	}
	return
}

func flush(c <-chan *ninchat.Event) {
	select {
	case _, open := <-c:
		if !open {
			return
		}

	default:
	}

	go func() {
		for range c {
		}
	}()
}
