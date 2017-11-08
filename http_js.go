// +build js

package ninchat

import (
	"errors"

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
		if xhrType == js.Undefined {
			xhrType = js.Module.Get("require").Invoke("xhr2")
		}
		xhrRequestHeaderSupport = true
	}
}

type httpHeader map[string]string

func (h httpHeader) Set(key, value string) {
	h[key] = value
}

type httpRequest struct {
	Method string
	URL    string
	Header httpHeader

	data *js.Object
}

func newGETRequest(url string) (req *httpRequest, err error) {
	req = &httpRequest{
		Method: "GET",
		URL:    url,
		Header: make(httpHeader),
	}
	return
}

func newDataRequest(method, url string, data *js.Object) (req *httpRequest, err error) {
	req = &httpRequest{
		Method: method,
		URL:    url,
		Header: make(httpHeader),
		data:   data,
	}
	return
}

type httpResponse struct {
	data *js.Object
	err  error
}

func getResponseData(req *httpRequest, timeout duration) (*js.Object, error) {
	resp := <-getResponseChannel(req, timeout)
	return resp.data, resp.err
}

// putResponseChannel makes an XMLHttpRequest, or an XDomainRequest in case of
// IE 8 or 9.
func putResponseToChannel(req *httpRequest, timeout duration, c chan<- httpResponse) {
	defer func() {
		if err := jsError(recover()); err != nil {
			c <- httpResponse{err: err}
		}
	}()

	xhr := xhrType.New()

	xhr.Set("onload", func() {
		response := xhr.Get("responseText")
		go func() {
			c <- httpResponse{data: response}
		}()
	})

	xhr.Set("onprogress", func() {
		// https://stackoverflow.com/questions/7037627/long-poll-and-ies-xdomainrequest-object
		js.Global.Call("setTimeout", func() {}, 0)
	})

	xhr.Set("ontimeout", func() {
		go func() {
			c <- httpResponse{err: errors.New("timeout")}
		}()
	})

	xhr.Set("onerror", func() {
		go func() {
			c <- httpResponse{err: errors.New("error")}
		}()
	})

	xhr.Call("open", req.Method, req.URL)
	xhr.Set("timeout", timeout)

	if xhrRequestHeaderSupport {
		for key, value := range req.Header {
			xhr.Call("setRequestHeader", key, value)
		}
	}

	js.Global.Call("setTimeout", func() {
		// http://cypressnorth.com/programming/internet-explorer-aborting-ajax-requests-fixed/
		xhr.Call("send", req.data)
	}, 0)
}
