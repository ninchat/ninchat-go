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

	// assume that JSONP works always (endpoint discovery worked, after all)
	connWorked = true

	url := "https://" + host + pollPath

	if s.sessionId == nil {
		s.log("session creation")

		header := s.makeCreateSessionAction()

		response, err := DataJSONP(url, header, JitterDuration(sessionCreateTimeout, 0.2))
		if err != nil {
			s.log("session creation:", err)
			return
		}

		select {
		case array := <-response:
			if array == nil {
				s.log("session creation timeout")
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

		gotOnline = true
	} else {
		s.log("session resumption")
	}

	if longPollTransfer(s, url) {
		gotOnline = true
	}

	return
}

// longPollTransfer sends buffered actions and polls for events.  It stops if
// two subsequent requests (of any kind) time out.
func longPollTransfer(s *Session, url string) (gotOnline bool) {
	var poller <-chan js.Object
	var sender <-chan js.Object
	var sendingId uint64
	var timeouts int

	for timeouts < 2 {
		if poller == nil {
			var err error

			header := s.makeResumeSessionAction(true)

			poller, err = DataJSONP(url, header, JitterDuration(minPollTimeout, 0.2))
			if err != nil {
				s.log("poll:", err)
				return
			}
		}

		if sender == nil && s.numSent < len(s.sendBuffer) {
			action := s.sendBuffer[s.numSent]

			if action.Payload != nil {
				json := action.Payload.Index(0).Str()

				object, err := ParseJSON(json)
				if err != nil {
					s.log("send:", err)
					return
				}

				action.Header.Set("payload", object)
			}

			action.Header.Set("session_id", s.sessionId)

			channel, err := DataJSONP(url, action.Header, JitterDuration(minSendTimeout, 0.2))

			action.Header.Delete("payload")
			action.Header.Delete("session_id")

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

		var array js.Object

		select {
		case array = <-poller:
			if array == nil {
				s.log("poll timeout")
			}

			poller = nil

		case array = <-sender:
			if array == nil {
				s.log("send timeout")
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

		case <-s.closeNotify:
			longPollClose(s, url)
			return
		}

		if array == nil {
			timeouts++
			continue
		}

		timeouts = 0

		for i := 0; i < array.Length(); i++ {
			header := array.Index(i)
			payload := NewArray()

			if object := header.Get("payload"); !object.IsUndefined() {
				json, err := StringifyJSON(object)
				if err != nil {
					s.log("poll payload:", err)
					return
				}

				payload.Call("push", json)
			}

			ackedActionId, _, ok := s.handleEvent(header, payload)

			// poll acked the action being sent before we got send response?
			if sendingId > 0 && sendingId <= ackedActionId {
				sendingId = 0
				s.numSent++
			}

			if !ok {
				return
			}

			gotOnline = true
		}
	}

	return
}

// longPollClose sends a close_session action without caring about the
// response.
func longPollClose(s *Session, url string) {
	header := map[string]interface{}{
		"action":     "close_session",
		"session_id": s.sessionId,
	}

	if _, err := DataJSONP(url, header, JitterDuration(minSendTimeout, 0.9)); err != nil {
		s.log("send:", err)
	}
}
