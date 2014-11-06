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

func PostCall(header js.Object, log func(...interface{}), address string) (channel chan js.Object, err error) {
	defer func() {
		err = jsError(recover())
	}()

	json, err := StringifyJSON(header)
	if err != nil {
		return
	}

	url := "https://" + address + callPath

	channel = make(chan js.Object, 1)

	xhr := js.Global.Get("XMLHttpRequest").New()

	xhr.Set("onload", func() {
		defer func() {
			if x := recover(); x != nil {
				println(x)
			}
		}()

		var array js.Object

		if xhr.Get("status").Int() == 200 {
			object, err := ParseJSON(xhr.Get("response").Str())
			if err != nil {
				log(err)
			} else {
				array = NewArray()
				array.Call("push", object)
			}
		} else {
			log("call status", xhr.Get("status").Str())
		}

		go func() {
			channel <- array
		}()
	})

	xhr.Call("open", "POST", url, true)
	xhr.Call("setRequestHeader", "Content-Type", "application/json")
	xhr.Call("send", json)

	return
}
