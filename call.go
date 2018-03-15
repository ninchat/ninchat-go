package ninchat

import (
	"errors"
)

const (
	callTimeout = second * 11
)

var defaultCaller Caller

// Call the sessionless API.
//
// If the action's OnReply callback is set, it will be called with the result
// events, unless an error occurs.
//
// https://ninchat.com/api#sessionless-http-calling
//
func Call(action *Action) (events []*Event, err error) {
	return defaultCaller.Call(action)
}

type Caller struct {
	Address string
}

func (caller *Caller) Call(action *Action) (events []*Event, err error) {
	url := "https://" + getAddress(caller.Address) + callPath

	req, err := newJSONRequest(url, action.Params)
	if err != nil {
		return
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	timeout := jitterDuration(callTimeout, 0.1)

	data, err := getResponseData(req, timeout)
	if err != nil {
		return
	}

	var headers []interface{}

	if jsonUnmarshalArray(data, &headers) != nil {
		var params map[string]interface{}

		if err = jsonUnmarshalObject(data, &params); err != nil {
			return
		}

		headers = []interface{}{params}
	}

	events = make([]*Event, 0, len(headers))

	for i, xParams := range headers {
		params, ok := xParams.(map[string]interface{})
		if !ok {
			err = errors.New("response event header is not an object")
			return
		}

		event := &Event{
			Params:    params,
			LastReply: i == len(headers)-1,
		}

		if action.OnReply != nil {
			action.OnReply(event)
		}

		events = append(events, event)
	}

	return
}
