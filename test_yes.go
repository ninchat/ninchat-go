// +build test

package ninchat

var defaultAddress = "api.ninchat.com"

type testSupport struct {
	ws           *webSocket
	forgetting   bool
	disconnected chan struct{}
	reconnect    chan struct{}

	forgotten []*webSocket
}

func (t *testSupport) setWebSocket(ws *webSocket) {
	t.ws = ws
}

func (t *testSupport) shouldAbortConn() bool {
	return t.forgetting
}

func (s *Session) closeWebSocket(ws *webSocket) {
	if s.test.forgetting {
		s.test.forgotten = append(s.test.forgotten, ws)
	} else {
		ws.close()
	}
}

func (t *testSupport) webSocketDisconnected() {
	t.ws = nil
	t.forgetting = false

	if t.disconnected == nil {
		return
	}

	close(t.disconnected)
	t.disconnected = nil

	if t.reconnect == nil {
		return
	}

	<-t.reconnect
}

func (s *Session) TestBreakConnection() (reconnector func()) {
	return s.testProblem(s.testDisconnect)
}

func (s *Session) TestForgetConnection() (reconnector func()) {
	return s.testProblem(func() bool {
		s.test.forgetting = true
		return true
	})
}

func (s *Session) testProblem(op func() bool) (reconnector func()) {
	if s.test.ws == nil {
		return
	}

	dis := make(chan struct{})
	s.test.disconnected = dis

	if !op() {
		return
	}

	re := make(chan struct{})
	s.test.reconnect = re

	<-dis

	return func() {
		close(re)
	}
}

func (s *Session) TestLoseSession() {
	s.sessionId = "lost"
}
