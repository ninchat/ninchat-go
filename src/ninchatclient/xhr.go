package main

import (
	"github.com/gopherjs/gopherjs/js"
)

var (
	xhrType                 *js.Object
	xhrRequestHeaderSupport bool
)

func init() {
	xhrType = js.Global.Get("XDomainRequest")
	if xhrType == js.Undefined {
		xhrType = js.Global.Get("XMLHttpRequest")
		xhrRequestHeaderSupport = true
	}
}

// XHR makes an XMLHttpRequest, or an XDomainRequest in case of IE 8 or 9.
// channel produces a response content string on success, an empty string on
// error, or nothing on timeout.
func XHR(url string, data string, timeout Duration) (channel chan string, err error) {
	defer func() {
		err = jsError(recover())
	}()

	var method string
	if data == "" {
		method = "GET"
	} else {
		method = "POST"
	}

	channel = make(chan string, 1)

	request := xhrType.New()

	request.Set("onload", func() {
		var response string

		if obj := request.Get("responseText"); obj != js.Undefined && obj != nil {
			response = obj.String()
		}

		go func() {
			channel <- response
		}()
	})

	request.Set("onprogress", func() {
		// https://stackoverflow.com/questions/7037627/long-poll-and-ies-xdomainrequest-object
	})

	request.Set("ontimeout", func() {
		go func() {
			close(channel)
		}()
	})

	request.Set("onerror", func() {
		go func() {
			channel <- ""
		}()
	})

	request.Call("open", method, url)
	request.Set("timeout", timeout)

	if data != "" && xhrRequestHeaderSupport {
		request.Call("setRequestHeader", "Content-Type", "application/json")
	}

	request.Call("send", data)

	return
}

// XHR_JSON stringifies data and calls XHR.
func XHR_JSON(url string, data interface{}, timeout Duration) (channel chan string, err error) {
	json, err := StringifyJSON(data)
	if err != nil {
		return
	}

	return XHR(url, json, timeout)
}
