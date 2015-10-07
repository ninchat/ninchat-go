package ninchatmodel

import (
	"errors"

	"github.com/ninchat/ninchat-go"
	api "github.com/ninchat/ninchat-go/ninchatapi"
)

// State
type State struct {
	ninchat.Session
	Self      SelfState
	Settings  SettingsState
	Users     UserState
	Dialogues DialogueState

	OnSessionEvent func(*ninchat.Event)
	OnEvent        func(*ninchat.Event)
}

// Open the session.  Client code must have called SetParams and initialized
// the compulsory callbacks...
func (state *State) Open() (err error) {
	state.Self.init(&state.Session)
	state.Settings.init(&state.Session)
	state.Users.init(&state.Session)
	state.Dialogues.init(&state.Session)

	initChan := make(chan error, 1)
	initDone := false

	state.Session.OnSessionEvent = func(e *ninchat.Event) {
		err := state.handle(e)

		if !initDone {
			select {
			case initChan <- err:
			default:
			}

			initDone = true
		} else if err != nil {
			state.log("session:", err)
		}

		if state.OnSessionEvent != nil {
			state.OnSessionEvent(e)
		}
	}

	state.Session.OnEvent = func(e *ninchat.Event) {
		if err := state.handle(e); err != nil {
			state.log(e.String(), "event:", err)
		}

		if state.OnEvent != nil {
			state.OnEvent(e)
		}
	}

	state.Session.Open()

	err = <-initChan
	if err != nil {
		state.Close()
	}
	return
}

func (state *State) handle(clientEvent *ninchat.Event) (err error) {
	defer func() {
		if x := recover(); x != nil {
			state.log("panic during event handling:", x)
			err = errors.New("event handler panicked")
		}
	}()

	event, err := api.NewEvent(clientEvent)
	if err != nil {
		return
	}

	switch e := event.(type) {
	case *api.Error:
		err = e

	case *api.SessionCreated:
		initial := (state.Self.Id == "")
		if initial {
			state.Users.Map[e.UserId] = &state.Self.User
		}

		state.Self.handleSession(e)
		state.Settings.handle(e.UserSettings, e.String())

		for peerId := range e.UserDialogues {
			state.Users.discover(peerId, true)
		}
		state.Dialogues.handleUser(state.Self.Id, e.UserDialogues, e.String())

		if !initial {
			for _, d := range state.Dialogues.Map {
				d.Window.handleSecondarySession(&state.Session, d)
			}
		}

	case *api.SessionStatusUpdated:
		if e.UserId != nil {
			state.Dialogues.handleSessionStatus(e)
		}

	case *api.UserFound:
		if e.UserId == state.Self.Id {
			state.Self.handleUser(e.UserAttrs, e.String())
			state.Settings.handle(e.UserSettings, e.String())

			for peerId := range e.UserDialogues {
				state.Users.discover(peerId, false)
			}
			state.Dialogues.handleUser(state.Self.Id, e.UserDialogues, e.String())
		} else {
			state.Users.handle(e.UserId, e.UserAttrs, e.String())
		}

	case *api.UserUpdated:
		if e.UserId == state.Self.Id {
			state.Self.handleUser(e.UserAttrs, e.String())
			state.Settings.handle(e.UserSettings, e.String())
		} else {
			state.Users.handle(e.UserId, e.UserAttrs, e.String())
		}

	case *api.DialogueUpdated:
		if e.DialogueStatus != nil {
			state.Users.discover(e.UserId, false)
		}
		state.Dialogues.handleDialogue(state.Self.Id, e)

	case *api.MessageReceived:
		if e.MessageUserId != nil {
			state.Users.discover(*e.MessageUserId, false)
		}

		if e.UserId != nil {
			state.Users.discover(*e.UserId, false)
			state.Dialogues.handleReceive(state.Self.Id, e)
		}

	case *api.MessageUpdated:
		// TODO: state.Channels.handleUpdate(state.Self.Id, e)
	}

	return
}

func (state *State) log(fragments ...interface{}) {
	if state.OnLog != nil {
		state.OnLog(fragments...)
	}
}

func log(session *ninchat.Session, prefix string, fragments []interface{}) {
	if session.OnLog != nil {
		fragments = append([]interface{}{prefix}, fragments...)
		session.OnLog(fragments...)
	}
}
