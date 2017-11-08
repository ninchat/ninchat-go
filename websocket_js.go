// +build js

package ninchat

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
)

var (
	webSocketClass     *js.Object
	webSocketSupported bool
)

func init() {
	// https://stackoverflow.com/questions/13349305/web-sockets-on-samsung-galaxy-s3-android-browser
	webSocketClass = js.Global.Get("WebSocket")
	if webSocketClass == js.Undefined {
		webSocketClass = js.Module.Get("require").Invoke("ws")
	}
	webSocketSupported = (webSocketClass != js.Undefined && webSocketClass.Get("CLOSING") != js.Undefined)
}

type webSocket struct {
	notify    chan struct{}
	goingAway bool
	err       error

	impl *js.Object
	open bool
	buf  []*js.Object
}

func newWebSocket(url string, timeout duration) (ws *webSocket) {
	ws = &webSocket{
		notify: make(chan struct{}, 1),
		impl:   webSocketClass.New(url),
	}

	ws.impl.Set("binaryType", "arraybuffer")

	var notifyClosed bool

	closeNotify := func() {
		if !notifyClosed {
			notifyClosed = true

			go func() {
				close(ws.notify)
			}()
		}
	}

	timeoutId := js.Global.Call("setTimeout", func() {
		ws.err = errors.New("timeout")
		closeNotify()
	}, timeout)

	clearTimeout := func() {
		if timeoutId != nil {
			js.Global.Call("clearTimeout", timeoutId)
			timeoutId = nil
		}
	}

	ws.impl.Set("onopen", func(*js.Object) {
		clearTimeout()

		if ws.err != nil {
			ws.impl.Call("close")
			return
		}

		ws.open = true

		go func() {
			if !notifyClosed {
				ws.notify <- struct{}{}
			}
		}()
	})

	ws.impl.Set("onmessage", func(object *js.Object) {
		ws.buf = append(ws.buf, object.Get("data"))

		go func() {
			if !notifyClosed {
				select {
				case ws.notify <- struct{}{}:
				default:
				}
			}
		}()
	})

	ws.impl.Set("onclose", func(object *js.Object) {
		ws.goingAway = (object.Get("code").Int() == 1001)
		ws.open = false
		closeNotify()
	})

	ws.impl.Set("onerror", func(object *js.Object) {
		clearTimeout()

		if ws.err == nil {
			ws.err = errors.New(object.Get("message").String())
		}

		closeNotify()
	})

	return
}

func (ws *webSocket) send(data *js.Object) (err error) {
	defer func() {
		err = jsError(recover())
	}()

	ws.impl.Call("send", data)
	return
}

func (ws *webSocket) sendJSON(object map[string]interface{}) (err error) {
	json, err := jsonMarshal(object)
	if err != nil {
		return
	}

	err = ws.send(json)
	return
}

func (ws *webSocket) sendPayload(action *Action) (err error) {
	defer func() {
		err = jsError(recover())
	}()

	decodeDataURI := (action.String() == "update_user")

	for _, frame := range action.Payload {
		var data *js.Object = frame

		if decodeDataURI {
			base64 := data.Call("split", ",").Index(1)
			binaryString := js.Global.Call("atob", base64)
			length := binaryString.Length()
			data = js.Global.Get("ArrayBuffer").New(length)
			array := js.Global.Get("Uint8Array").New(data)

			for i := 0; i < length; i++ {
				array.SetIndex(i, binaryString.Call("charCodeAt", i))
			}
		}

		ws.impl.Call("send", data)
	}

	return
}

func (ws *webSocket) receive() (x *js.Object) {
	if len(ws.buf) > 0 {
		x = ws.buf[0]
		ws.buf = ws.buf[1:]
	}
	return
}

func (ws *webSocket) receiveJSON() (object map[string]interface{}, err error) {
	x := ws.receive()
	if x == nil {
		return
	}

	err = jsonUnmarshalObject(stringData(x), &object)
	return
}

func (ws *webSocket) close() {
	defer func() {
		recover()
	}()

	ws.impl.Call("close")
}
