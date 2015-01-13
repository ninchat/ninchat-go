package main

import (
	"github.com/gopherjs/gopherjs/js"
)

const (
	minPollTimeout = Second * 64
	minSendTimeout = Second * 7
)

// LongPollTransport creates or resumes a session, and runs an I/O loop after
// that.
func LongPollTransport(s *Session, host string) (connWorked, gotOnline bool) {
	defer func() {
		if err := jsError(recover()); err != nil {
			s.log("poll:", err)
		}
	}()

	url := "https://" + host + pollPath

	if s.sessionId == nil {
		s.log("session creation")

		header := s.makeCreateSessionAction()

		creator, err := XHR_JSON(url, header, JitterDuration(sessionCreateTimeout, 0.2))
		if err != nil {
			s.log("session creation:", err)
			return
		}

		select {
		case response, ok := <-creator:
			if !ok {
				s.log("session creation timeout")
				return
			} else if response == "" {
				s.log("session creation error")
				return
			}

			array, err := ParseJSON(response)
			if err != nil {
				s.log("session creation response:", err)
				return
			}

			header := array.Index(0)

			if !s.handleSessionEvent(header) {
				return
			}

		case <-s.closeNotify:
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
		if err := longPollPing(s, url); err != nil {
			s.log("session resumption:", err)
			return
		}
	}

	longPollTransfer(s, url, &connWorked, &gotOnline)
	return
}

// longPollTransfer sends buffered actions and polls for events.  It stops if
// two subsequent requests (of any kind) time out.
func longPollTransfer(s *Session, url string, connWorked, gotOnline *bool) {
	var poller <-chan string
	var sender <-chan string
	var sendingId uint64
	var failures int

	s.numSent = 0

	for failures < 2 {
		if poller == nil {
			var err error

			header := s.makeResumeSessionAction(true)

			poller, err = XHR_JSON(url, header, JitterDuration(minPollTimeout, 0.2))
			if err != nil {
				s.log("poll:", err)
				return
			}
		}

		if sender == nil && s.numSent < len(s.sendBuffer) {
			action := s.sendBuffer[s.numSent]

			if action.Payload != nil {
				var payload js.Object
				var err error

				frame := action.Payload.Index(0)

				if action.Header.Get("action").Str() == "update_user" {
					base64, err := ParseDataURI(frame)
					if err != nil {
						s.log("send:", err)
						return
					}

					payload = NewArray()
					payload.Call("push", base64)
				} else {
					if payload, err = ParseJSON(frame.Str()); err != nil {
						s.log("send:", err)
						return
					}
				}

				action.Header.Set("payload", payload)
			}

			action.Header.Set("session_id", s.sessionId)

			request, err := StringifyJSON(action.Header)

			action.Header.Delete("session_id")
			action.Header.Delete("payload")

			if err != nil {
				s.log("send:", err)
				return
			}

			channel, err := XHR(url, request, JitterDuration(minSendTimeout, 0.2))
			if err != nil {
				s.log("send:", err)
				return
			}

			if action.Id == 0 {
				s.sendBuffer = append(s.sendBuffer[:s.numSent], s.sendBuffer[s.numSent+1:]...)
			} else {
				sender = channel
				sendingId = action.Id
			}
		}

		var response string
		var ok bool

		select {
		case response, ok = <-poller:
			if !ok {
				s.log("poll timeout")
			} else if response == "" {
				s.log("poll error")
			}

			poller = nil
			s.connActive()

		case response, ok = <-sender:
			if !ok {
				s.log("send timeout")
			} else if response == "" {
				s.log("send error")
			} else if sendingId > 0 {
				s.numSent++
			}

			sender = nil
			sendingId = 0

		case sending := <-s.sendNotify:
			if !sending {
				longPollClose(s, url)
				return
			}

			continue

		case <-s.closeNotify:
			longPollClose(s, url)
			return
		}

		var array js.Object

		if response != "" {
			var err error

			array, err = ParseJSON(response)
			if err != nil {
				s.log("response:", err)
			}
		}

		if array == nil {
			failures++
			s.numSent = 0
			continue
		}

		failures = 0
		*connWorked = true

		for i := 0; i < array.Length(); i++ {
			header := array.Index(i)
			payload := NewArray()

			if object := header.Get("payload"); object != js.Undefined {
				json, err := StringifyJSON(object)
				if err != nil {
					s.log("poll payload:", err)
					return
				}

				payload.Call("push", json)
			}

			ackedActionId, sessionLost, _, ok := s.handleEvent(header, payload)

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
func longPollPing(s *Session, url string) (err error) {
	header := map[string]interface{}{
		"action":     "ping",
		"session_id": s.sessionId,
	}

	_, err = XHR_JSON(url, header, JitterDuration(minSendTimeout, 0.9))
	return
}

// longPollClose sends a close_session action without caring about the
// response.
func longPollClose(s *Session, url string) {
	header := map[string]interface{}{
		"action":     "close_session",
		"session_id": s.sessionId,
	}

	if _, err := XHR_JSON(url, header, JitterDuration(minSendTimeout, 0.9)); err != nil {
		s.log("send:", err)
	}
}
