package main

import (
	"github.com/gopherjs/gopherjs/js"
)

const (
	callTimeout = Second * 11

	callLogInvocationName = namespace + ".call onLog callback"
)

// Call implements the call(object[, function|null[, string|null]]) JavaScript
// API.
func Call(header, onLog, address js.Object) (promise interface{}) {
	url := "https://" + GetAddress(address) + callPath

	deferred, promise := Defer()

	go func() {
		channel, err := XHR_JSON(url, header, JitterDuration(callTimeout, 0.1))
		if err != nil {
			Log(callLogInvocationName, onLog, "call:", err)
			deferred.Reject()
			return
		}

		response, ok := <-channel
		if response == "" {
			if ok {
				Log(callLogInvocationName, onLog, "call error")
			} else {
				Log(callLogInvocationName, onLog, "call timeout")
			}

			deferred.Reject()
			return
		}

		event, err := ParseJSON(response)
		if err != nil {
			Log(callLogInvocationName, onLog, "call response:", err)
			deferred.Reject()
			return
		}

		events := NewArray()
		events.Call("push", event)

		deferred.Resolve(events)
	}()

	return
}
