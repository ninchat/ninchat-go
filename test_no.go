// +build !test

package ninchat

const defaultAddress = "api.ninchat.com"

type testSupport struct{}

func (*testSupport) setWebSocket(*webSocket) {}
func (*testSupport) webSocketDisconnected()  {}
func (*testSupport) shouldAbortConn() bool   { return false }

func (*Session) closeWebSocket(ws *webSocket) { ws.close() }
