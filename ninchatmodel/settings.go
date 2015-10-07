package ninchatmodel

import (
	"reflect"

	"github.com/ninchat/ninchat-go"
)

// SettingsState
type SettingsState struct {
	OnChange func(Change, map[string]interface{})
	Data     map[string]interface{}

	session *ninchat.Session
}

func (state *SettingsState) init(session *ninchat.Session) {
	state.session = session
}

func (state *SettingsState) handle(settings map[string]interface{}, eventName string) {
	if settings == nil {
		return
	}

	var c Change

	if state.Data != nil {
		if !reflect.DeepEqual(state.Data, settings) {
			state.Data = settings
			c = Updated
		}
	} else {
		state.Data = settings
		c = Added
	}

	state.log(c.String(), "by", eventName)

	if c != unchanged {
		state.OnChange(c, settings)
	}
}

func (state *SettingsState) log(fragments ...interface{}) {
	log(state.session, "settings:", fragments)
}
