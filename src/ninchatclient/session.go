package main

import (
	"sort"

	"github.com/gopherjs/gopherjs/js"
)

const (
	discoveryTimeout  = Second * 7
	connectTimeout    = Second * 9
	connectIterations = 2
	maxBackoffDelay   = Second * 60

	sessionCreateTimeout = Second * 13

	maxEventAckDelay  = Second * 7
	maxEventAckWindow = 4096

	sessionSessionEventInvocationName = namespace + ".Session onSessionEvent callback"
	sessionEventInvocationName        = namespace + ".Session onEvent callback"
	sessionLogInvocationName          = namespace + ".Session onLog callback"
)

var (
	sessionEventAckWindow = JitterUint64(maxEventAckWindow, -0.25)
)

// Transport is an interface implemented by WebSocketTransport and
// LongPollTransport.
type Transport func(s *Session, host string) (connWorked, gotOnline bool)

// Session hides the details of Ninchat API connection management.
type Session struct {
	onSessionEvent js.Object
	onEvent        js.Object
	onLog          js.Object
	address        string
	forceLongPoll  bool
	sessionParams  js.Object
	sessionId      js.Object

	lastActionId    uint64
	sendNotify      chan bool
	sendBuffer      []*Action
	numSent         int
	sendEventAck    bool
	receivedEventId uint64
	ackedEventId    uint64

	closeNotify chan bool
	closed      bool
}

// NewSession implements the newSession() JavaScript API.
func NewSession() map[string]interface{} {
	s := Session{
		address: defaultAddress,
		closed:  true,
	}

	return map[string]interface{}{
		"onSessionEvent": s.OnSessionEvent,
		"onEvent":        s.OnEvent,
		"onLog":          s.OnLog,
		"setParams":      s.SetParams,
		"setTransport":   s.SetTransport,
		"setAddress":     s.SetAddress,
		"open":           s.Open,
		"close":          s.Close,
		"send":           s.Send,
	}
}

// OnSessionEvent implements the Session.onSessionEvent(function) JavaScript
// API.
func (s *Session) OnSessionEvent(callback js.Object) {
	s.onSessionEvent = callback
}

// OnEvent implements the Session.onEvent(function) JavaScript API.
func (s *Session) OnEvent(callback js.Object) {
	s.onEvent = callback
}

// OnLog implements the Session.onLog(function|null) JavaScript API.
func (s *Session) OnLog(callback js.Object) {
	if callback.IsNull() {
		callback = nil
	}

	s.onLog = callback
}

// SetParams implements the Session.setParams(object) JavaScript API.
func (s *Session) SetParams(params js.Object) {
	if params.Get("message_types").IsUndefined() {
		panic("message_types parameter not defined")
	}

	s.sessionParams = params
}

// SetTransport implements the Session.setTransport(string|null) JavaScript
// API.
func (s *Session) SetTransport(name js.Object) {
	if name.IsNull() {
		s.forceLongPoll = false
		return
	}

	switch string := name.Str(); string {
	case "websocket":
		panic("websocket transport cannot be forced")

	case "longpoll":
		s.forceLongPoll = true

	default:
		panic("unknown transport: " + string)
	}
}

// SetAddress implements the Session.setAddress(string|null) JavaScript API.
func (s *Session) SetAddress(address js.Object) {
	s.address = GetAddress(address)
}

// Open implements the Session.open() JavaScript API.
func (s *Session) Open() {
	if s.sendNotify != nil {
		panic("session already initialized")
	}

	if s.onSessionEvent == nil {
		panic("onSessionEvent callback not defined")
	}

	if s.onEvent == nil {
		panic("onEvent callback not defined")
	}

	if s.sessionParams == nil {
		panic("session parameters not defined")
	}

	s.sendNotify = make(chan bool, 1)
	s.closeNotify = make(chan bool, 1)
	s.closed = false

	go s.discover()
}

// Close implements the Session.close() JavaScript API.
func (s *Session) Close() {
	if s.closed {
		return
	}

	for _, action := range s.sendBuffer {
		s.resolve(action, false)
	}

	s.sendBuffer = nil
	s.numSent = 0
	s.closed = true

	go func() {
		s.closeNotify <- true
		close(s.sendNotify)
	}()
}

// Send implements the Session.send(object[, array]) JavaScript API.
func (s *Session) Send(header, payload js.Object) (promise js.Object) {
	if s.sendNotify == nil {
		panic("session not initialized")
	}

	if s.closed {
		panic("session already closed")
	}

	if payload.IsUndefined() || payload.IsNull() || payload.Length() == 0 {
		payload = nil
	}

	action := &Action{
		Header:  header,
		Payload: payload,
	}

	if header.Get("action_id").IsNull() {
		header.Delete("action_id")
	} else {
		s.lastActionId++
		action.Id = s.lastActionId
		header.Set("action_id", action.Id)

		promise, action.Resolve = NewPromise()
	}

	s.send(action)
	return
}

// send buffers an action.
func (s *Session) send(action *Action) {
	s.sendBuffer = append(s.sendBuffer, action)

	go func() {
		select {
		case s.sendNotify <- true:
		default:
		}
	}()

	return
}

// sendAck indicates that this is a good time to acknowledge received events.
func (s *Session) sendAck() {
	s.sendEventAck = true

	go func() {
		select {
		case s.sendNotify <- true:
		default:
		}
	}()
}

// reset clears session state, but preserves actions in the send buffer.
func (s *Session) reset() {
	s.sessionId = nil

	if len(s.sendBuffer) == 0 {
		s.lastActionId = 0
	}

	s.numSent = 0
	s.sendEventAck = false
	s.receivedEventId = 0
	s.ackedEventId = 0
}

// discover runs an endpoint discovery loop.
func (s *Session) discover() {
	s.log("opening")
	defer s.log("closed")

	var backoff Backoff
	var wsFailed bool

	for !s.closed {
		s.log("endpoint discovery")

		url := "https://" + s.address + endpointPath

		if response, err := JSONP(url, JitterDuration(discoveryTimeout, 0.1)); err != nil {
			s.log("endpoint discovery:", err)
		} else {
			select {
			case endpoint := <-response:
				if endpoint == nil {
					s.log("endpoint discovery timed out")
				} else if hosts, err := GetEndpointHosts(endpoint); err != nil {
					s.log("endpoint discovery:", err)
				} else {
					s.log("endpoint discovered")

					if !WebSocketSupported || s.forceLongPoll || wsFailed {
						s.connect(LongPollTransport, hosts, &backoff)
					} else {
						wsFailed = !s.connect(WebSocketTransport, hosts, &backoff)
					}

					// no delay; discovery was successful even if connect wasn't
					continue
				}

			case <-s.closeNotify:
				return
			}
		}

		if delay := backoff.Failure(maxBackoffDelay); delay > 0 {
			s.log("sleeping")
			Sleep(delay)
		}
	}
}

// connect tries to connect to each host a few times using the given transport
// implementation.
func (s *Session) connect(transport Transport, hosts []string, backoff *Backoff) (transportWorked bool) {
	for trial := 0; trial < connectIterations; trial++ {
		for _, host := range hosts {
			//gopherjs:blocking
			connWorked, gotOnline := transport(s, host)

			if connWorked {
				transportWorked = true
			}

			if gotOnline {
				backoff.Success()
				return
			}

			if s.closed {
				return
			}

			if delay := backoff.Failure(maxBackoffDelay); delay > 0 {
				s.log("sleeping")
				Sleep(delay)
			}
		}
	}

	return
}

// makeCreateSessionAction makes a create_session action header.
func (s *Session) makeCreateSessionAction() (header js.Object) {
	header = s.sessionParams
	header.Set("action", "create_session")

	return
}

// makeResumeSessionAction makes a resume_session action header for
// initializing new connections, or polling and acknowledging events.
func (s *Session) makeResumeSessionAction(session bool) (header js.Object) {
	header = NewObject()
	header.Set("action", "resume_session")
	header.Set("event_id", s.receivedEventId)

	if session {
		header.Set("session_id", s.sessionId)
	}

	s.sendEventAck = false
	s.ackedEventId = s.receivedEventId

	return
}

// handleSessionEvent establishes a session or fails.
func (s *Session) handleSessionEvent(header js.Object) (ok bool) {
	userId, userAuth, sessionId, eventId, ok, err := GetSessionEventCredentials(header)

	if err != nil {
		s.log("session creation:", err)
	}

	if !jsInvoke(sessionSessionEventInvocationName, s.onSessionEvent, header) {
		ok = false
	}

	if !ok {
		s.Close()
		return
	}

	s.sessionParams.Set("user_id", userId)

	if userAuth != nil {
		s.sessionParams.Set("user_auth", userAuth)
	}

	for _, param := range []string{"identity_type", "identity_name", "identity_auth"} {
		if newValue := s.sessionParams.Get(param + "_new"); !newValue.IsUndefined() {
			s.sessionParams.Set(param, newValue)
		}
	}

	s.sessionId = sessionId
	s.receivedEventId = eventId

	s.log("session created")

	ok = true
	return
}

// handleEvent parses the event header partially, does whatever is necessary,
// and usually passes the event to client code.
func (s *Session) handleEvent(header, payload js.Object) (actionId uint64, needsAck, ok bool) {
	eventId, actionId, err := GetEventAndActionId(header)
	if err != nil {
		s.log("event:", err)
		return
	}

	if eventId > 0 {
		s.receivedEventId = eventId

		if !s.sendEventAck {
			if s.receivedEventId-s.ackedEventId >= sessionEventAckWindow {
				s.sendAck()
			} else {
				needsAck = true
			}
		}
	}

	if actionId > 0 {
		i := sort.Search(s.numSent, func(i int) bool {
			action := s.sendBuffer[i]
			return action.Id >= actionId
		})

		if i < s.numSent {
			if action := s.sendBuffer[i]; action.Id == actionId {
				s.resolve(action, true, header, payload)

				lastReply, err := IsEventLastReply(header, action)
				if err != nil {
					s.log("event:", err)
					return
				}

				if lastReply {
					s.sendBuffer = append(s.sendBuffer[:i], s.sendBuffer[i+1:]...)
					s.numSent--
				}
			}
		}
	}

	errorType, errorReason, sessionLost, err := GetEventError(header)
	if err != nil {
		s.log("event:", err)
		if sessionLost {
			s.reset()
		}
		return
	}

	if errorType == "deprecated" {
		s.log("deprecated:", errorReason)
	}

	if !jsInvoke(sessionEventInvocationName, s.onEvent, header, payload) {
		return
	}

	ok = true
	return
}

// resolve fulfills the promise made by the Send method.
func (s *Session) resolve(action *Action, successful bool, args ...interface{}) {
	defer func() {
		if x := recover(); x != nil {
			println(x)
		}
	}()

	if action.Resolve != nil {
		action.Resolve(successful, args...)
	}
}

// log conditionally passes a message to the client code.
func (s *Session) log(tokens ...interface{}) {
	Log(sessionLogInvocationName, s.onLog, tokens...)
}
