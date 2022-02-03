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

func getUserAgent(custom map[string][]string) string {
	if vv, found := custom[userAgentHeader]; found && len(vv) > 0 {
		return vv[0]
	}

	return DefaultUserAgent
}

type webSocket struct {
	notify    chan struct{}
	goingAway bool
	err       error

	impl *js.Object
	open bool
	buf  []*js.Object

	client string
}

func newWebSocket(host string, header map[string][]string, timeout duration, log func(...interface{})) *webSocket {
	deadline := timeAdd(timeNow(), timeout)

	ws := &webSocket{
		notify: make(chan struct{}, 1),
		client: getUserAgent(header),
	}

	var (
		poked        bool
		notifyClosed bool
		connect      func() *js.Object
	)

	connect = func() *js.Object {
		impl := webSocketClass.New("wss://" + host + socketPath)
		impl.Set("binaryType", "arraybuffer")

		stale := false

		closeNotify := func() {
			if !notifyClosed {
				notifyClosed = true

				go func() {
					close(ws.notify)
				}()
			}
		}

		timeoutId := js.Global.Call("setTimeout", func() {
			if stale {
				return
			}

			ws.err = errors.New("timeout")
			closeNotify()
		}, timeout)

		clearTimeout := func() {
			if timeoutId != nil {
				js.Global.Call("clearTimeout", timeoutId)
				timeoutId = nil
			}
		}

		impl.Set("onopen", func(*js.Object) {
			clearTimeout()
			if stale {
				return
			}
			if ws.err != nil {
				impl.Call("close")
				return
			}

			ws.open = true

			go func() {
				if !notifyClosed {
					ws.notify <- struct{}{}
				}
			}()
		})

		impl.Set("onmessage", func(object *js.Object) {
			if stale {
				return
			}

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

		impl.Set("onclose", func(object *js.Object) {
			if stale {
				return
			}

			ws.goingAway = (object.Get("code").Int() == 1001)
			ws.open = false
			closeNotify()
		})

		impl.Set("onerror", func(object *js.Object) {
			clearTimeout()
			if stale {
				return
			}
			if ws.err != nil {
				closeNotify()
				return
			}

			var pendingErr error
			if msg := object.Get("message"); msg != js.Undefined {
				pendingErr = errors.New("websocket onerror: message: " + msg.String())
			} else {
				pendingErr = errors.New("websocket onerror: " + object.String())
			}

			timeout = timeSub(deadline, timeNow())
			if poked || timeout <= 0 {
				ws.err = pendingErr
				closeNotify()
				return
			}

			origin := "https://" + host
			log("poking", origin)

			req, err := newGetRequest(origin+endpointPath, nil)
			if err != nil {
				panic(err)
			}

			poked = true
			stale = true

			go func() {
				if _, err := getResponseData(req, timeout); err == nil {
					timeout = timeSub(deadline, timeNow())
					if timeout > 0 {
						log("retrying same websocket host once more")
						ws.impl = connect()
						return
					}
				} else {
					log("poke error:", err)
				}

				ws.err = pendingErr
				closeNotify()
			}()
		})

		return impl
	}

	ws.impl = connect()
	return ws
}

func (ws *webSocket) sendInitialJSON(object map[string]interface{}) error {
	if ws.client != "" {
		object["client"] = ws.client
	}
	return ws.sendJSON(object)
}

func (ws *webSocket) send(data *js.Object) (err error) {
	defer func() {
		err = jsError("websocket send", recover())
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
		err = jsError("websocket send payload", recover())
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
