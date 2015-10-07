// +build !js

package ninchat

import (
	"crypto/tls"
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const webSocketSupported = true

var webSocketTLSConfig = &tls.Config{
	PreferServerCipherSuites: true,
}

type webSocket struct {
	notify    chan struct{}
	goingAway bool
	err       error

	conn *websocket.Conn
	lock sync.Mutex
	buf  [][]byte
}

func newWebSocket(url string, timeout duration) (ws *webSocket) {
	dialer := &websocket.Dialer{
		TLSClientConfig:  webSocketTLSConfig,
		HandshakeTimeout: time.Duration(timeout),
	}

	ws = &webSocket{
		notify: make(chan struct{}, 1),
	}

	if ws.conn, _, ws.err = dialer.Dial(url, nil); ws.err != nil {
		close(ws.notify)
		return
	}

	ws.notify <- struct{}{}

	go func() {
		defer close(ws.notify)

		for {
			var (
				typ  int
				data []byte
			)

			if typ, data, ws.err = ws.conn.ReadMessage(); ws.err != nil {
				if err, ok := ws.err.(*websocket.CloseError); ok {
					ws.goingAway = (err.Code == websocket.CloseGoingAway)
				}
				return
			}

			switch typ {
			case websocket.TextMessage, websocket.BinaryMessage:
				ws.lock.Lock()
				ws.buf = append(ws.buf, data)
				ws.lock.Unlock()

				select {
				case ws.notify <- struct{}{}:
				default:
				}

			default:
				// XXX: is this possible?
			}
		}
	}()

	return
}

func (ws *webSocket) send(data []byte) error {
	return ws.conn.WriteMessage(websocket.BinaryMessage, data)
}

func (ws *webSocket) sendJSON(object interface{}) error {
	return ws.conn.WriteJSON(object)
}

func (ws *webSocket) sendPayload(action *Action) (err error) {
	for _, data := range action.Payload {
		if err = ws.send(data); err != nil {
			return
		}
	}
	return
}

func (ws *webSocket) receive() (data []byte) {
	ws.lock.Lock()
	defer ws.lock.Unlock()

	if len(ws.buf) > 0 {
		data = ws.buf[0]
		ws.buf = ws.buf[1:]
	}

	return
}

func (ws *webSocket) receiveJSON() (object map[string]interface{}, err error) {
	x := ws.receive()
	if x == nil {
		return
	}

	err = json.Unmarshal(x, &object)
	return
}

func (ws *webSocket) close() {
	if ws.conn != nil {
		ws.conn.Close()
	}
}
