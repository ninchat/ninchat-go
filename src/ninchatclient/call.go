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
func Call(header, onLog, address js.Object) (promise js.Object) {
	url := "https://" + GetAddress(address) + callPath

	promise, resolve := NewPromise()

	go func() {
		response, err := DataJSONP(url, header, JitterDuration(callTimeout, 0.1))
		if err != nil {
			Log(callLogInvocationName, onLog, "call:", err)
			resolve(false)
			return
		}

		if events := <-response; events != nil {
			resolve(true, events)
		} else {
			resolve(false)
		}
	}()

	return
}
