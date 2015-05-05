package main

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
)

var (
	WebSocketSupported = js.Global.Get("WebSocket") != js.Undefined
)

// WebSocket
type WebSocket struct {
	Notify    chan bool
	GoingAway bool

	impl   *js.Object
	open   bool
	error  error
	buffer []*js.Object
}

// NewWebSocket
func NewWebSocket(url string) (ws *WebSocket) {
	ws = &WebSocket{
		Notify: make(chan bool, 1),
		impl:   js.Global.Get("WebSocket").New(url),
	}

	ws.impl.Set("binaryType", "arraybuffer")

	ws.impl.Set("onopen", func(*js.Object) {
		ws.open = true

		go func() {
			ws.Notify <- true
		}()
	})

	ws.impl.Set("onmessage", func(object *js.Object) {
		ws.buffer = append(ws.buffer, object.Get("data"))

		go func() {
			select {
			case ws.Notify <- true:
			default:
			}
		}()
	})

	ws.impl.Set("onclose", func(object *js.Object) {
		ws.GoingAway = object.Get("code").Int() == 1001
		ws.open = false

		go func() {
			close(ws.Notify)
		}()
	})

	ws.impl.Set("onerror", func(object *js.Object) {
		ws.error = errors.New("WebSocket error event")
	})

	return
}

// Send
func (ws *WebSocket) Send(data interface{}) (err error) {
	defer func() {
		err = jsError(recover())
	}()

	if err = ws.error; err != nil {
		return
	}

	ws.impl.Call("send", data)
	return
}

// SendJSON
func (ws *WebSocket) SendJSON(object interface{}) (err error) {
	json, err := StringifyJSON(object)
	if err != nil {
		return
	}

	err = ws.Send(json)
	return
}

// Receive
func (ws *WebSocket) Receive() (data *js.Object, err error) {
	if err = ws.error; err != nil {
		return
	}

	if !ws.open {
		return
	}

	if len(ws.buffer) > 0 {
		data = ws.buffer[0]
		ws.buffer = ws.buffer[1:]
	}

	return
}

// ReceiveJSON
func (ws *WebSocket) ReceiveJSON() (object *js.Object, err error) {
	data, err := ws.Receive()
	if err != nil || data == nil {
		return
	}

	return ParseJSON(StringifyFrame(data))
}

// Close
func (ws *WebSocket) Close() (err error) {
	defer func() {
		err = jsError(recover())
	}()

	ws.impl.Call("close")

	err = ws.error
	return
}

// StringifyFrame
func StringifyFrame(data *js.Object) (s string) {
	s, ok := data.Interface().(string)
	if ok {
		return
	}

	view := NewUint8Array(data)
	bytes := view.Interface().([]uint8)
	s = string(bytes)

	return
}
