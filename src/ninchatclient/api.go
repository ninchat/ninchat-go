package main

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
)

const (
	defaultAddress = "api.ninchat.com"
	protocolPath   = "/v2"
	endpointPath   = protocolPath + "/endpoint"
	socketPath     = protocolPath + "/socket"
	pollPath       = protocolPath + "/poll"
	callPath       = protocolPath + "/call"
)

// Action
type Action struct {
	Id       uint64
	Header   js.Object
	Payload  js.Object
	Deferred *Deferred

	name string
}

func (a *Action) Name() string {
	if a.name == "" {
		a.name = a.Header.Get("action").Str()
	}

	return a.name
}

// GetAddress
func GetAddress(address js.Object) string {
	if address == js.Undefined || address == nil {
		return defaultAddress
	} else {
		return address.Str()
	}
}

// GetEndpointHosts
func GetEndpointHosts(response string) (hosts []string, err error) {
	defer func() {
		if e := jsError(recover()); e != nil {
			err = e
		}
	}()

	endpoint, err := ParseJSON(response)
	if err != nil {
		return
	}

	jsHosts := endpoint.Get("hosts")

	if jsHosts.Length() == 0 {
		err = errors.New("endpoint hosts array is empty")
		return
	}

	hosts = make([]string, jsHosts.Length())

	for i := 0; i < jsHosts.Length(); i++ {
		hosts[i] = jsHosts.Index(i).Str()
	}

	return
}

// GetSessionEventCredentials
func GetSessionEventCredentials(header js.Object) (userId, userAuth, sessionId js.Object, eventId uint64, ok bool, err error) {
	defer func() {
		if e := jsError(recover()); e != nil {
			err = e
		}
	}()

	if header.Get("event").Str() != "session_created" {
		return
	}

	userId = header.Get("user_id")

	if object := header.Get("user_auth"); object != js.Undefined {
		userAuth = object
	}

	sessionId = header.Get("session_id")
	eventId = header.Get("event_id").Uint64()

	ok = true
	return
}

// GetEventFrames
func GetEventFrames(header js.Object) (frames int, err error) {
	defer func() {
		err = jsError(recover())
	}()

	if object := header.Get("frames"); object != js.Undefined {
		if frames = object.Int(); frames < 0 {
			frames = 0
		}
	}

	return
}

// GetEventAndActionId
func GetEventAndActionId(header js.Object) (eventId uint64, actionId uint64, err error) {
	defer func() {
		err = jsError(recover())
	}()

	if object := header.Get("event_id"); object != js.Undefined {
		eventId = object.Uint64()
	}

	if object := header.Get("action_id"); object != js.Undefined {
		actionId = object.Uint64()
	}

	return
}

// IsEventLastReply
func IsEventLastReply(header js.Object, action *Action) (lastReply bool, err error) {
	defer func() {
		err = jsError(recover())
	}()

	lastReply = true

	if historyLength := header.Get("history_length"); historyLength != js.Undefined {
		if historyLength.Int() > 0 {
			lastReply = false
		}
	}

	if action.name == "search" {
		users := header.Get("users")
		channels := header.Get("channels")

		if users != js.Undefined || channels != js.Undefined {
			lastReply = false
		}
	}

	return
}

// GetEventError
func GetEventError(header js.Object) (errorType, errorReason string, sessionLost bool, err error) {
	defer func() {
		if e := jsError(recover()); e != nil {
			err = e
		}
	}()

	if header.Get("event").Str() != "error" {
		return
	}

	errorType = header.Get("error_type").Str()

	if object := header.Get("error_reason"); object != js.Undefined {
		errorReason = object.Str()
	}

	switch errorType {
	case "session_not_found":
		sessionLost = true
		fallthrough

	case "connection_superseded", "message_has_too_many_parts", "message_part_too_long", "message_too_long", "request_malformed":
		if errorReason != "" {
			err = errors.New("error: " + errorType + " (" + errorReason + ")")
		} else {
			err = errors.New("error: " + errorType)
		}
	}

	return
}
