package ninchatmodel

import (
	"github.com/ninchat/ninchat-go"
	api "github.com/ninchat/ninchat-go/ninchatapi"
)

// SelfState
type SelfState struct {
	User
	OnChange func(c Change, u *User, auth string)

	session *ninchat.Session
}

func (state *SelfState) init(session *ninchat.Session) {
	state.session = session
}

func (state *SelfState) handleSession(e *api.SessionCreated) {
	var (
		c    Change
		auth string
	)

	if state.Id != "" {
		c = state.update(e.UserAttrs)
	} else {
		state.Id = e.UserId
		state.Attrs = e.UserAttrs
		if e.UserAuth != nil {
			auth = *e.UserAuth
		}
		c = Added
	}

	state.log(state.Id, c.String(), "by", e.String())

	if c != unchanged {
		state.OnChange(c, &state.User, auth)
	}
}

func (state *SelfState) handleUser(attrs *api.UserAttrs, eventName string) {
	c := state.update(attrs)

	state.log(state.Id, c.String(), "by", eventName)

	if c != unchanged {
		state.OnChange(c, &state.User, "")
	}
}

func (state *SelfState) log(fragments ...interface{}) {
	log(state.session, "self:", fragments)
}
