// Ninchat API connection library.
//
// https://ninchat.com/api
//
package ninchat

import (
	"sort"
	"sync"
)

const (
	discoveryTimeout  = second * 7
	connectTimeout    = second * 9
	connectIterations = 2
	maxBackoffDelay   = second * 60

	sessionCreateTimeout = second * 13

	maxEventAckDelay  = second * 7
	maxEventAckWindow = 4096
)

var (
	sessionEventAckWindow = jitterInt64(maxEventAckWindow, -0.25)
)

// Session hides the details of API connection management.  It needs to be
// initialized by setting at least OnSessionEvent and OnEvent, and calling
// SetParams.  After that the Open method is used to make a connection to the
// server.  Finally, the Close method disconnects from the server.
type Session struct {
	// OnSessionEvent is the session creation handler.  It will be invoked with
	// a "session_created" or an "error" event.
	//
	// If another "session_created" event is received, it means that the previous
	// session was lost, and a new one was established automatically.
	//
	// If an "error" event is received, it means that a new session can't be
	// established without intervention.  The client code must call SetParams to
	// supply new credentials, unless it decides to Close.
	//
	OnSessionEvent func(*Event)

	// OnEvent is the handler for in-session events.
	//
	// "error" events received via this callback are not fatal.
	//
	OnEvent func(*Event)

	// OnClose is an optional session closure handler.  It will be invoked
	// after a Close call has been fully processed.  It won't be invoked if an
	// "error" event is received via OnSessionEvent (unless SetParams is called
	// again).
	//
	OnClose func()

	// OnConnState is an optional connection state change monitor.  It will be
	// called with one of the following strings:
	//
	// - "connecting"
	// - "connected"
	// - "disconnected"
	//
	OnConnState func(state string)

	// OnConnActive is an optional connection activity monitor.  It will be
	// called (approximately) when data is received on the underlying network
	// connection.
	OnConnActive func()

	// OnLog is an optional message logger.
	OnLog func(fragments ...interface{})

	Address string

	forceLongPoll bool // only for testing

	mutex sync.Mutex // guards all variables below

	sessionParams map[string]interface{}
	sessionId     interface{}

	latestConnState string

	lastActionId    int64
	sendNotify      chan struct{}
	sendBuffer      []*Action
	numSent         int
	sendEventAck    bool
	receivedEventId int64
	ackedEventId    int64

	closeNotify chan struct{}
	closed      bool
	running     bool
}

// transport is an interface implemented by webSocketTransport and
// longPollTransport.
type transport func(s *Session, host string) (connWorked, gotOnline bool)

// SetParams sets "create_session" action parameters.  If Open has already been
// called, this takes effect when a session is lost.
func (s *Session) SetParams(params map[string]interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if params["message_types"] == nil {
		panic("message_types parameter not defined")
	}

	if x, found := params["session_id"]; x != nil {
		s.sessionId = x
	} else if found {
		delete(params, "session_id")
	}

	s.sessionParams = params

	if s.sendNotify != nil && !s.running {
		// opened + stopped -> restart
		s.running = true
		go s.discover()
	}
}

// Open creates a session on the server.
func (s *Session) Open() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.closed {
		panic("session already closed")
	}

	if s.sendNotify != nil {
		panic("session already initialized")
	}

	if s.OnSessionEvent == nil {
		panic("onSessionEvent callback not defined")
	}

	if s.OnEvent == nil {
		panic("onEvent callback not defined")
	}

	if s.sessionParams == nil {
		panic("session parameters not defined")
	}

	s.sendNotify = make(chan struct{}, 1)
	s.closeNotify = make(chan struct{}, 1)
	s.running = true

	go s.discover()
}

// Close the session on the server.
func (s *Session) Close() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.closed {
		return
	}

	for _, action := range s.sendBuffer {
		if action.OnReply != nil {
			action.OnReply(nil)
		}
	}

	s.sendBuffer = nil
	s.numSent = 0

	s.closed = true
	s.running = false

	go func() {
		s.closeNotify <- struct{}{}
		close(s.sendNotify)
	}()
}

// Send an action.
//
// To send an action without an "action_id" parameter, specify it as nil.
// Otherwise an "action_id" is generated automatically.
//
// If "action_id" is used and the action's OnReply callback is set, it will be
// called when a reply event is received.  If the Session object is closed
// before a reply is received, OnReply will be called with a nil argument.
//
// With specific actions that cause multiple reply events, the OnReply callback
// will be called for each event until the final event.  Only the last event
// will have the LastReply member set.
//
func (s *Session) Send(action *Action) (err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.sendNotify == nil {
		panic("session not initialized")
	}

	if s.closed {
		panic("session already closed")
	}

	if len(action.Payload) == 0 {
		action.Payload = nil
	}

	if x, found := action.Params["action_id"]; found && x == nil {
		delete(action.Params, "action_id")
	} else {
		s.lastActionId++
		action.id = s.lastActionId
		action.Params["action_id"] = action.id
	}

	s.send(action)
	return
}

// send buffers an action.
func (s *Session) send(action *Action) {
	s.sendBuffer = append(s.sendBuffer, action)

	go func() {
		select {
		case s.sendNotify <- struct{}{}:
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
		case s.sendNotify <- struct{}{}:
		default:
		}
	}()
}

// discover runs an endpoint discovery loop.
func (s *Session) discover() {
	s.mutex.Lock()
	defer func() {
		closed := s.closed
		s.mutex.Unlock()
		if closed && s.OnClose != nil {
			s.OnClose()
		}
	}()

	s.log("starting")
	defer s.log("stopped")

	defer s.connState("disconnected")

	var backoff backoff

	for s.running {
		s.connState("connecting")

		url := "https://" + getAddress(s.Address) + endpointPath

		request, err := newGETRequest(url)
		if err != nil {
			panic(err)
		}

		s.mutex.Unlock()
		select {
		case response := <-getResponseChannel(request, jitterDuration(discoveryTimeout, 0.1)):
			s.mutex.Lock()

			var hosts []string
			err := response.err

			if err == nil {
				var endpoint map[string]interface{}

				if err = jsonUnmarshalObject(response.data, &endpoint); err == nil {
					hosts, err = getEndpointHosts(endpoint)
				}
			}

			if err == nil {
				s.log("endpoint discovered")

				if webSocketSupported && !s.forceLongPoll {
					if s.connect(webSocketTransport, hosts, &backoff) {
						continue
					}
				}

				s.connect(longPollTransport, hosts, &backoff)
			} else {
				s.log("endpoint discovery:", err)
			}

		case <-s.closeNotify:
			s.mutex.Lock()
			return
		}

		if !s.backOff(&backoff) {
			return
		}
	}
}

// connect tries to connect to each host a few times using the given transport
// implementation.
func (s *Session) connect(transport transport, hosts []string, backoff *backoff) (transportWorked bool) {
	for trial := 0; trial < connectIterations; trial++ {
		for _, host := range hosts {
			s.connState("connecting")

			for i, c := range s.Address {
				if c == '/' {
					host += s.Address[i:]
					break
				}
			}

			//gopherjs:blocking
			connWorked, gotOnline := transport(s, host)

			if connWorked {
				transportWorked = true
			}

			if gotOnline {
				backoff.success()
				return
			}

			if !s.running {
				return
			}

			if !s.backOff(backoff) {
				return
			}
		}
	}

	return
}

// backOff sleeps for a time.  False is returned if session was closed while
// sleeping.
func (s *Session) backOff(b *backoff) (ok bool) {
	delay := b.failure(maxBackoffDelay)
	if delay == 0 {
		ok = true
		return
	}

	s.connState("disconnected")
	s.log("sleeping")

	s.mutex.Unlock()
	defer s.mutex.Lock()

	select {
	case <-newTimer(delay).C:
		ok = true

	case <-s.closeNotify:
	}

	return
}

// canLogin checks whether the makeCreateSessionAction method would make an
// action that has any chance of creating a session for an existing user.
// The answer is no if the action could succeed, but would create a new
// user.
func (s *Session) canLogin() bool {
	if s.sessionParams["access_key"] != nil {
		return true
	}

	if s.sessionParams["user_id"] != nil {
		return s.sessionParams["user_auth"] != nil || s.sessionParams["master_sign"] != nil
	}

	return s.sessionParams["identity_type"] != nil && s.sessionParams["identity_name"] != nil && s.sessionParams["identity_auth"] != nil
}

// makeCreateSessionAction makes a create_session action header.
func (s *Session) makeCreateSessionAction() (params map[string]interface{}) {
	params = map[string]interface{}{
		"action": "create_session",
	}

	if userId := s.sessionParams["user_id"]; userId == nil {
		// Client code is responsible for specifying correct parameters.

		for key, value := range s.sessionParams {
			params[key] = value
		}
	} else {
		// This might be automatic session recreation, try to be smart.

		if userAuth := s.sessionParams["user_auth"]; userAuth != nil {
			params["user_id"] = userId
			params["user_auth"] = userAuth
		} else if masterSign := s.sessionParams["master_sign"]; masterSign != nil {
			params["user_id"] = userId
			params["master_sign"] = masterSign
		} else if identityType := s.sessionParams["identity_type"]; identityType != nil {
			params["identity_type"] = identityType
			params["identity_name"] = s.sessionParams["identity_name"]
			params["identity_auth"] = s.sessionParams["identity_auth"]
		} else {
			// Fallback: let the server decide.  (But still make sure that we
			// won't be creating a new user by accident.)

			params["user_id"] = userId
		}

		for key, value := range s.sessionParams {
			switch key {
			case "user_id", "user_auth", "identity_type", "identity_name", "identity_auth", "access_key", "master_sign":
				// skipped

			default:
				params[key] = value
			}
		}
	}

	return
}

// makeResumeSessionAction makes a resume_session action header for
// initializing new connections, or polling and acknowledging events.
func (s *Session) makeResumeSessionAction(session bool) (params map[string]interface{}) {
	params = map[string]interface{}{
		"action":   "resume_session",
		"event_id": s.receivedEventId,
	}

	if session {
		params["session_id"] = s.sessionId
	}

	s.sendEventAck = false
	s.ackedEventId = s.receivedEventId

	return
}

// handleSessionEvent establishes a session or fails.
func (s *Session) handleSessionEvent(params map[string]interface{}) (ok bool) {
	event := &Event{
		Params: params,
	}

	quit := false

	if event.String() == "error" {
		s.sessionId = nil
		quit = true

		switch errorType, _ := event.Str("error_type"); errorType {
		case "internal":
			// keep trying

		default:
			s.running = false
		}
	}

	s.deliverSessionEvent(event)

	if quit {
		return
	}

	delete(s.sessionParams, "user_attrs")
	delete(s.sessionParams, "user_settings")
	delete(s.sessionParams, "identity_attrs")
	delete(s.sessionParams, "access_key")
	delete(s.sessionParams, "master_sign")

	s.sessionParams["user_id"] = event.Params["user_id"]

	if x := event.Params["user_auth"]; x != nil {
		s.sessionParams["user_auth"] = x
	}

	for _, param := range []string{"identity_type", "identity_name", "identity_auth"} {
		if newValue := s.sessionParams[param+"_new"]; newValue != nil {
			delete(s.sessionParams, param+"_new")
			s.sessionParams[param] = newValue
		}
	}

	s.sessionId = event.Params["session_id"]

	if len(s.sendBuffer) == 0 {
		s.lastActionId = 0
	}

	s.sendEventAck = false
	s.receivedEventId, _ = event.Int64("event_id")
	s.ackedEventId = 0

	s.log("session created")

	ok = true
	return
}

// handleEvent parses the event header partially, does whatever is necessary,
// and usually passes the event to client code.
func (s *Session) handleEvent(event *Event) (actionId int64, sessionLost, needsAck, ok bool) {
	if eventId, _ := event.Int64("event_id"); eventId > 0 {
		s.receivedEventId = eventId

		if !s.sendEventAck {
			if s.receivedEventId-s.ackedEventId >= sessionEventAckWindow {
				s.sendAck()
			} else {
				needsAck = true
			}
		}
	}

	actionId, _ = event.Int64("action_id")
	if actionId > 0 {
		i := sort.Search(s.numSent, func(i int) bool {
			action := s.sendBuffer[i]
			return action.id >= actionId
		})

		if i < s.numSent {
			if action := s.sendBuffer[i]; action.id == actionId {
				event.initLastReply(action)

				if action.OnReply != nil {
					action.OnReply(event)
				}

				if event.LastReply {
					s.sendBuffer = append(s.sendBuffer[:i], s.sendBuffer[i+1:]...)
					s.numSent--
				}
			}
		}
	}

	if event.String() == "user_deleted" {
		s.sessionId = nil
		s.running = false
		s.deliverSessionEvent(event)

		sessionLost = true
		return
	}

	errorType, errorReason, sessionLost, err := event.getError()
	if err != nil {
		s.log("event:", err)

		if sessionLost {
			s.sessionId = nil

			if !s.canLogin() {
				s.running = false
				s.deliverSessionEvent(event)
			}
		}

		return
	}

	if errorType == "deprecated" {
		s.log("deprecated:", errorReason)
	}

	s.deliverEvent(event)

	ok = true
	return
}

func (s *Session) deliverSessionEvent(event *Event) {
	s.mutex.Unlock()
	defer s.mutex.Lock()

	s.OnSessionEvent(event)
}

func (s *Session) deliverEvent(event *Event) {
	s.mutex.Unlock()
	defer s.mutex.Lock()

	s.OnEvent(event)
}

// connState passes an enumeration value to the client code.
func (s *Session) connState(state string) {
	if s.latestConnState == state {
		return
	}

	s.latestConnState = state

	if s.OnConnState == nil {
		return
	}

	s.mutex.Unlock()
	defer s.mutex.Lock()

	s.OnConnState(state)
}

// connActive pokes the client code.
func (s *Session) connActive() {
	if s.OnConnActive == nil {
		return
	}

	s.mutex.Unlock()
	defer s.mutex.Lock()

	s.OnConnActive()
}

// log passes a message to the client code.
func (s *Session) log(tokens ...interface{}) {
	if s.OnLog == nil {
		return
	}

	s.mutex.Unlock()
	defer s.mutex.Lock()

	s.OnLog(tokens...)
}
