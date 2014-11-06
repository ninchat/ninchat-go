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
	if address.IsUndefined() || address.IsNull() {
		return defaultAddress
	} else {
		return address.Str()
	}
}

// GetEndpointHosts
func GetEndpointHosts(endpoint js.Object) (hosts []string, err error) {
	defer func() {
		if e := jsError(recover()); e != nil {
			err = e
		}
	}()

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

	if object := header.Get("user_auth"); !object.IsUndefined() {
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

	if object := header.Get("frames"); !object.IsUndefined() {
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

	if object := header.Get("event_id"); !object.IsUndefined() {
		eventId = object.Uint64()
	}

	if object := header.Get("action_id"); !object.IsUndefined() {
		actionId = object.Uint64()
	}

	return
}

// IsEventLastReply
func IsEventLastReply(header js.Object, action *Action) (lastReply bool, err error) {
	defer func() {
		err = jsError(recover())
	}()

	switch action.name {
	case "load_history":
		historyLength := header.Get("history_length").Int()
		lastReply = (historyLength == 0)

	case "search":
		users := header.Get("users")
		channels := header.Get("channels")
		lastReply = users.IsUndefined() && channels.IsUndefined()

	default:
		lastReply = true
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

	if object := header.Get("error_reason"); !object.IsUndefined() {
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
