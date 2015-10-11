package ninchatapi

import (
	"github.com/ninchat/ninchat-go"
)

// Action interface is implemented by all action structs.
type Action interface {
	newClientAction() (*ninchat.Action, error)
}

// Call an action with or without a session.
func Call(session *ninchat.Session, events chan<- *ninchat.Event, action Action) (err error) {
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

	if session == nil {
		if _, err = ninchat.Call(clientAction); err != nil {
			close(events)
			return
		}
	} else {
		if x, found := clientAction.Params["action_id"]; found && x == nil {
			close(events)
			panic("calling via session but action_id is disabled")
		}

		session.Send(clientAction)
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

func unaryCall(session *ninchat.Session, action Action, event Event) (ok bool, err error) {
	c := make(chan *ninchat.Event, 1) // XXX: why doesn't this work without buffering?

	if err = Call(session, c, action); err != nil {
		return
	}

	clientEvent := <-c
	flush(c)

	if clientEvent == nil {
		return
	}

	ok = true

	if clientEvent.String() == "error" {
		err = newError(clientEvent)
	} else {
		err = event.MergeFrom(clientEvent)
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
