package ninchat

import (
	"errors"
)

const (
	defaultAddress = "api.ninchat.com"
	protocolPath   = "/v2"
	endpointPath   = protocolPath + "/endpoint"
	socketPath     = protocolPath + "/socket"
	pollPath       = protocolPath + "/poll"
	callPath       = protocolPath + "/call"
)

var (
	errNoActionId = errors.New("no action_id")
)

// getAddress
func getAddress(address string) string {
	if address == "" {
		return defaultAddress
	} else {
		return address
	}
}

// getEndpointHosts
func getEndpointHosts(object map[string]interface{}) (hosts []string, err error) {
	x := object["hosts"]
	if x == nil {
		err = errors.New("invalid endpoint document")
		return
	}

	xHosts, ok := x.([]interface{})
	if !ok {
		err = errors.New("invalid endpoint hosts type")
		return
	}

	if len(xHosts) == 0 {
		err = errors.New("no endpoint hosts")
		return
	}

	for _, x := range xHosts {
		if s, ok := x.(string); ok {
			hosts = append(hosts, s)
		} else {
			err = errors.New("invalid endpoint host value type")
		}
	}

	if len(hosts) > 0 {
		err = nil
	}

	return
}

// String returns the action's name.
func (a *Action) String() (s string) {
	if x, ok := a.Params["action"]; ok {
		s, _ = x.(string)
	}
	return
}

// GetId returns the action_id if Send has been called and it generated one.
func (a *Action) GetId() (id int64, err error) {
	id = a.id
	if id == 0 {
		err = errNoActionId
	}
	return
}

// Bool looks up a boolean parameter.  A boolean parameter is considered false
// when it is nil or not set.
func (e *Event) Bool(param string) (v bool) {
	if x := e.Params[param]; x != nil {
		v, _ = x.(bool)
	}
	return
}

// Int looks up a numeric parameter and converts it to an integer.
func (e *Event) Int(param string) (v int, ok bool) {
	if x := e.Params[param]; x != nil {
		var f float64
		if f, ok = x.(float64); ok {
			v = int(f)
		}
	}
	return
}

// Int64 looks up a numeric parameter and converts it to a 64-bit integer.
func (e *Event) Int64(param string) (v int64, ok bool) {
	if x := e.Params[param]; x != nil {
		var f float64
		if f, ok = x.(float64); ok {
			v = int64(f)
		}
	}
	return
}

// Float64 looks up a numeric parameter.
func (e *Event) Float64(param string) (v float64, ok bool) {
	if x := e.Params[param]; x != nil {
		v, ok = x.(float64)
	}
	return
}

// Str looks up a string parameter.
func (e *Event) Str(param string) (v string, ok bool) {
	if x := e.Params[param]; x != nil {
		v, ok = x.(string)
	}
	return
}

// Array looks up an array parameter.
func (e *Event) Array(param string) (v []interface{}, ok bool) {
	if x := e.Params[param]; x != nil {
		v, ok = x.([]interface{})
	}
	return
}

// Map looks up an object parameter.
func (e *Event) Map(param string) (v map[string]interface{}, ok bool) {
	if x := e.Params[param]; x != nil {
		v, ok = x.(map[string]interface{})
	}
	return
}

// String returns the event's name.
func (e *Event) String() (s string) {
	s, _ = e.Str("event")
	return
}

func (e *Event) initLastReply(action *Action) {
	if n, _ := e.Int("history_length"); n > 0 {
		return
	}

	if e.String() == "search_results" {
		if e.Params["users"] != nil || e.Params["channels"] != nil {
			return
		}
	}

	e.LastReply = true
}

func (event *Event) getError() (errorType, errorReason string, sessionLost bool, err error) {
	if x, found := event.Params["event"]; !found || x.(string) != "error" {
		return
	}

	errorType = event.Params["error_type"].(string)

	if x := event.Params["error_reason"]; x != nil {
		errorReason = x.(string)
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
