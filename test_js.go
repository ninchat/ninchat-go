// +build test,js

package ninchat

func (s *Session) testDisconnect() (ok bool) {
	defer func() {
		if x := recover(); x != nil {
			msg := "panic"
			if err, ok := x.(error); ok {
				msg = err.Error()
			}
			println("TestDisconnect:", msg)
		}
	}()

	s.test.ws.impl.Call("close", 4321)
	ok = true
	return
}
