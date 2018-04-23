package ninchat

const (
	maxKeepaliveInterval = second * 56
	minReceiveTimeout    = second * 64
)

// webSocketTransport runs a reconnect loop for a given host.  It doesn't mind
// disconnections, but a failure to connect and establish a session stops the
// loop.
func webSocketTransport(s *Session, host string) (connWorked, gotOnline bool) {
	var ws *webSocket

	defer func() {
		if ws != nil {
			ws.close()
		}
	}()

	for {
		gotOnline = false
		hostHealthy := false

		s.log("connecting to", host)

		ws = newWebSocket("wss://"+host+socketPath, jitterDuration(connectTimeout, 0.1))

		s.mutex.Unlock()
		select {
		case _, connected := <-ws.notify:
			s.mutex.Lock()

			if connected {
				s.connState("connected")

				connWorked = true
				gotOnline, hostHealthy = webSocketHandshake(s, ws)
			} else if ws.err != nil {
				s.log("connection failed:", ws.err)
			} else {
				s.log("connection failed")
			}

		case <-s.closeNotify:
			s.mutex.Lock()
		}

		goingAway := ws.goingAway

		ws.close()
		ws = nil

		if goingAway {
			s.log("disconnected (server is going away)")
		} else {
			s.log("disconnected")
		}

		if !gotOnline || !hostHealthy || !s.running || goingAway {
			return
		}
	}
}

// webSocketHandshake creates or resumes a session, and runs I/O loops after
// that.
func webSocketHandshake(s *Session, ws *webSocket) (gotOnline, hostHealthy bool) {
	var params map[string]interface{}

	if s.sessionId == nil {
		s.log("session creation")
		params = s.makeCreateSessionAction()
	} else {
		s.log("session resumption")
		params = s.makeResumeSessionAction(true)
	}

	if err := ws.sendJSON(params); err != nil {
		s.log("send:", err)
	}

	if s.sessionId == nil {
		var params map[string]interface{}

		connected := true
		timer := newTimer(jitterDuration(sessionCreateTimeout, 0.2))

		for {
			var err error

			if params, err = ws.receiveJSON(); err != nil {
				s.log("session creation:", err)
				return
			}
			if params != nil {
				timer.Stop()
				break
			}

			if !connected {
				s.log("disconnected during session creation:", ws.err)
				timer.Stop()
				return
			}

			s.mutex.Unlock()
			select {
			case _, connected = <-ws.notify:
				s.mutex.Lock()

			case <-timer.C:
				s.mutex.Lock()

				s.log("session creation timeout")
				return
			}
		}

		if !s.handleSessionEvent(params) {
			return
		}

		gotOnline = true
		hostHealthy = true
		s.connActive()
	}

	fail := make(chan struct{}, 2)
	done := make(chan struct{})

	go func() {
		defer close(done)

		s.mutex.Lock()
		defer s.mutex.Unlock()

		webSocketSend(s, ws, fail)
	}()

	gotEvents, hostHealthy := webSocketReceive(s, ws, fail)
	if gotEvents {
		gotOnline = true
	}

	s.mutex.Unlock()
	defer s.mutex.Lock()

	<-done
	return
}

// webSocketSend sends buffered actions and acknowledges events received by
// webSocketReceive.  It makes sure that something is sent to the server every
// now and then.
func webSocketSend(s *Session, ws *webSocket, fail chan struct{}) {
	keeper := newTimer(jitterDuration(maxKeepaliveInterval, -0.3))
	defer keeper.Stop()

	s.numSent = 0

	for {
		for s.numSent < len(s.sendBuffer) {
			action := s.sendBuffer[s.numSent]

			if action.Payload != nil {
				action.Params["frames"] = len(action.Payload)
			}

			if s.receivedEventId != s.ackedEventId {
				action.Params["event_id"] = s.receivedEventId

				s.sendEventAck = false
				s.ackedEventId = s.receivedEventId
			}

			err := ws.sendJSON(action.Params)

			delete(action.Params, "frames")
			delete(action.Params, "event_id")

			if err != nil {
				s.log("send:", err)
				fail <- struct{}{}
				return
			}

			if action.Payload != nil {
				if err := ws.sendPayload(action); err != nil {
					s.log("send:", err)
					fail <- struct{}{}
					return
				}
			}

			if action.id == 0 {
				s.sendBuffer = append(s.sendBuffer[:s.numSent], s.sendBuffer[s.numSent+1:]...)
			} else {
				s.numSent++
			}

			keeper.Reset(jitterDuration(maxKeepaliveInterval, -0.3))
		}

		if s.sendEventAck && s.receivedEventId != s.ackedEventId {
			action := s.makeResumeSessionAction(false)

			if err := ws.sendJSON(action); err != nil {
				s.log("send:", err)
				fail <- struct{}{}
				return
			}
		}

		s.mutex.Unlock()
		select {
		case _, sending := <-s.sendNotify:
			s.mutex.Lock()

			if !sending {
				closeSession := map[string]interface{}{
					"action": "close_session",
				}

				if err := ws.sendJSON(closeSession); err != nil {
					s.log("send:", err)
				}

				return
			}

		case <-keeper.C:
			s.mutex.Lock()

			// empty keepalive frame between actions
			if err := ws.send(emptyData()); err != nil {
				s.log("send:", err)
				fail <- struct{}{}
				return
			}

			keeper.Reset(jitterDuration(maxKeepaliveInterval, -0.3))

		case <-fail:
			s.mutex.Lock()
			return
		}
	}
}

// webSocketReceive receives events.  It stops if the server doesn't send
// anything for some time.
func webSocketReceive(s *Session, ws *webSocket, fail chan struct{}) (gotEvents, hostHealthy bool) {
	wsNotify := ws.notify
	connected := true

	var event *Event
	var frames int

	d := jitterDuration(minReceiveTimeout, 0.3)
	watchdogTime := timeAdd(timeNow(), d)
	watchdog := newTimer(d)
	defer watchdog.Stop()

	acker := newTimer(-1)
	defer acker.Stop()

	for {
		var ackNeeded bool

		for {
			if event == nil {
				data := ws.receive()
				if data == nil {
					break
				}

				d := jitterDuration(minReceiveTimeout, 0.7)
				watchdogTime = timeAdd(timeNow(), d)
				if !watchdog.Active() {
					watchdog.Reset(d)
				}
				s.connActive()

				text := stringData(data)
				if dataLength(text) == 0 {
					// keepalive frame
					continue
				}

				var params map[string]interface{}

				if err := jsonUnmarshalObject(text, &params); err != nil {
					s.log("receive:", err)
					hostHealthy = false
					fail <- struct{}{}
					return
				}

				event = &Event{
					Params: params,
				}

				if n, _ := event.Int("frames"); n > 0 {
					frames = n
				}
			} else {
				data := ws.receive()
				if data == nil {
					if !connected {
						// disconnected before complete payload was sent
						event = nil
						frames = 0
					}
					break
				}

				event.Payload = append(event.Payload, data)
				frames--
			}

			if frames == 0 {
				_, sessionLost, needsAck, ok := s.handleEvent(event)
				if !ok {
					if sessionLost {
						gotEvents = true
					} else {
						hostHealthy = false
					}

					fail <- struct{}{}
					return
				}

				if needsAck {
					ackNeeded = true
				}

				event = nil
				frames = 0

				gotEvents = true
				hostHealthy = true
			}

			s.mutex.Unlock()
			select {
			case <-s.closeNotify:
				s.mutex.Lock()
				return

			case <-fail:
				s.mutex.Lock()
				return

			default:
				s.mutex.Lock()
				// don't wait
			}
		}

		if !connected && event == nil {
			if ws.err != nil {
				s.log("receive:", ws.err)
			} else {
				s.log("receive disconnect")
			}
			fail <- struct{}{}
			return
		}

		if ackNeeded && !acker.Active() {
			acker.Reset(jitterDuration(maxEventAckDelay, -0.3))
		}

		s.mutex.Unlock()
		select {
		case _, connected = <-wsNotify:
			s.mutex.Lock()

			if !connected {
				wsNotify = nil
			}

		case <-watchdog.C:
			s.mutex.Lock()

			if remain := timeSub(watchdogTime, timeNow()); remain > millisecond*16 {
				watchdog.Reset(remain)
			} else {
				s.log("receive timeout")
				fail <- struct{}{}
				return
			}

		case <-acker.C:
			s.mutex.Lock()

			if !s.sendEventAck && s.ackedEventId != s.receivedEventId {
				s.sendAck()
			}

		case <-s.closeNotify:
			s.mutex.Lock()
			return

		case <-fail:
			s.mutex.Lock()
			return
		}
	}
}
