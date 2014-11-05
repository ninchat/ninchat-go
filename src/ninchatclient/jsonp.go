package main

import (
	"github.com/gopherjs/gopherjs/js"
)

const (
	callbackPrefix = "_callback"
)

var (
	callbackNum int
)

// JSONP
func JSONP(url string, timeout Duration) (channel chan js.Object, err error) {
	return doJSONP(url+"?callback=", timeout)
}

// DataJSONP
func DataJSONP(url string, data interface{}, timeout Duration) (channel chan js.Object, err error) {
	json, err := StringifyJSON(data)
	if err != nil {
		return
	}

	return doJSONP(url+"?data="+EncodeURIComponent(json)+"&callback=", timeout)
}

func doJSONP(url string, timeout Duration) (channel chan js.Object, err error) {
	defer func() {
		err = jsError(recover())
	}()

	callbackId, _ := StringifyJSON(callbackNum)
	callbackNum = (callbackNum + 1) & 0x7fffffff

	function := callbackPrefix + callbackId
	callback := namespace + "." + function

	channel = make(chan js.Object, 1)

	timeoutId := SetTimeout(func() {
		module.Delete(function)

		go func() {
			close(channel)
		}()
	}, timeout)

	module.Set(function, func(object js.Object) {
		ClearTimeout(timeoutId)
		module.Delete(function)

		go func() {
			channel <- object
		}()
	})

	document := js.Global.Get("document")

	script := document.Call("createElement", "script")
	script.Set("type", "text/javascript")
	script.Set("async", true)
	script.Set("src", url+callback)

	head := document.Call("getElementsByTagName", "head").Index(0)
	head.Call("appendChild", script)

	return
}
