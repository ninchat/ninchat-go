package ninchat

const (
	minPollTimeout = second * 64
	minSendTimeout = second * 7
)

// longPollTransport creates or resumes a session, and runs an I/O loop after
// that.
func longPollTransport(s *Session, host string) (connWorked, gotOnline bool) {
	url := "https://" + host + pollPath

	if s.sessionId == nil {
		s.log("session creation")

		action := s.makeCreateSessionAction()
		timeout := jitterDuration(sessionCreateTimeout, 0.2)

		s.mutex.Unlock()
		select {
		case response := <-getJSONRequestResponseChannel(url, action, timeout):
			s.mutex.Lock()

			if response.err != nil {
				s.log("session creation:", response.err)
				return
			}

			var array []interface{}

			if err := jsonUnmarshalArray(response.data, &array); err != nil {
				s.log("session creation response:", err)
				return
			}

			if len(array) == 0 {
				s.log("session creation response JSON array is empty")
				return
			}

			event, ok := array[0].(map[string]interface{})
			if !ok {
				s.log("session creation response header is not a JSON object")
				return
			}

			if !s.handleSessionEvent(event) {
				return
			}

		case <-s.closeNotify:
			s.mutex.Lock()

			longPollClose(s, url)
			return
		}

		connWorked = true
		gotOnline = true
		s.connState("connected")
		s.connActive()
	} else {
		s.log("session resumption")

		// this ensures that we get an event quickly if the connection works,
		// and can update the status
		longPollPing(s, url)
	}

	longPollTransfer(s, url, &connWorked, &gotOnline)
	return
}

// longPollTransfer sends buffered actions and polls for events.  It stops if
// two subsequent requests (of any kind) time out.
func longPollTransfer(s *Session, url string, connWorked, gotOnline *bool) {
	var poller <-chan httpResponse
	var sender <-chan httpResponse
	var sendingId int64
	var failures int

	s.numSent = 0

	for failures < 2 {
		if poller == nil {
			action := s.makeResumeSessionAction(true)
			timeout := jitterDuration(minPollTimeout, 0.2)
			poller = getJSONRequestResponseChannel(url, action, timeout)
		}

		if sender == nil && s.numSent < len(s.sendBuffer) {
			action := s.sendBuffer[s.numSent]

			if action.Payload != nil {
				if action.String() == "update_user" {
					action.Params["payload"] = longPollBinaryPayload(action)
				} else {
					var object map[string]interface{}

					if err := jsonUnmarshalObject(action.Payload[0], &object); err != nil {
						s.log("send:", err)
						return
					}

					action.Params["payload"] = object
				}
			}

			action.Params["session_id"] = s.sessionId

			timeout := jitterDuration(minSendTimeout, 0.2)
			channel := getJSONRequestResponseChannel(url, action.Params, timeout)

			delete(action.Params, "session_id")
			delete(action.Params, "payload")

			if action.id == 0 {
				go logErrorResponseLockless(s, channel, "send error:")
				s.sendBuffer = append(s.sendBuffer[:s.numSent], s.sendBuffer[s.numSent+1:]...)
			} else {
				sender = channel
				sendingId = action.id
			}
		}

		var response httpResponse

		s.mutex.Unlock()
		select {
		case response = <-poller:
			s.mutex.Lock()

			if response.err != nil {
				s.log("poll error:", response.err)
			}

			poller = nil
			s.connActive()

		case response = <-sender:
			s.mutex.Lock()

			if response.err != nil {
				s.log("send error:", response.err)
			} else if sendingId > 0 {
				s.numSent++
			}

			sender = nil
			sendingId = 0

		case _, sending := <-s.sendNotify:
			s.mutex.Lock()

			if !sending {
				longPollClose(s, url)
				return
			}

			continue

		case <-s.closeNotify:
			s.mutex.Lock()

			longPollClose(s, url)
			return
		}

		var array []interface{}

		err := response.err
		if err == nil {
			if err = jsonUnmarshalArray(response.data, &array); err != nil {
				s.log("response:", err)
			}
		}
		if err != nil {
			failures++
			s.numSent = 0
			continue
		}

		failures = 0
		*connWorked = true

		for _, x := range array {
			params, ok := x.(map[string]interface{})
			if !ok {
				s.log("poll event is not an object")
				return
			}

			event := &Event{
				Params: params,
			}

			if x := params["payload"]; x != nil {
				object, ok := x.(map[string]interface{})
				if !ok {
					s.log("poll payload is not an object")
					return
				}

				json, err := jsonMarshal(object)
				if err != nil {
					s.log("poll payload:", err)
					return
				}

				event.Payload = singleFrame(json)
			}

			ackedActionId, sessionLost, _, ok := s.handleEvent(event)

			// poll acked the action being sent before we got send response?
			if sendingId > 0 && sendingId <= ackedActionId {
				sendingId = 0
				s.numSent++
			}

			if !ok {
				if sessionLost {
					*gotOnline = true
				}

				return
			}

			if !*gotOnline {
				*gotOnline = true
				s.connState("connected")
			}
		}
	}

	return
}

// longPollPing sends a ping action without an action_id.
func longPollPing(s *Session, url string) {
	action := map[string]interface{}{
		"action":     "ping",
		"session_id": s.sessionId,
	}

	c := getJSONRequestResponseChannel(url, action, jitterDuration(minSendTimeout, 0.9))
	go logErrorResponseLockless(s, c, "ping error:")
}

// longPollClose sends a close_session action without caring about the
// response.
func longPollClose(s *Session, url string) {
	action := map[string]interface{}{
		"action":     "close_session",
		"session_id": s.sessionId,
	}

	c := getJSONRequestResponseChannel(url, action, jitterDuration(minSendTimeout, 0.9))
	go logErrorResponseLockless(s, c, "send error:")
}

func logErrorResponseLockless(s *Session, channel <-chan httpResponse, prefix string) {
	resp := <-channel
	if resp.err != nil && s.OnLog != nil {
		s.OnLog(prefix, resp.err)
	}
}
