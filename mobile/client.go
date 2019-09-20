package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"

	ninchat "github.com/ninchat/ninchat-go"
	"github.com/ninchat/ninchat-go/ninchatstate"
)

func asError(x interface{}) error {
	if x == nil {
		return nil
	}

	if err, ok := x.(error); ok {
		return err
	} else {
		return fmt.Errorf("%v", x)
	}
}

type JSON struct {
	x json.RawMessage
}

func NewJSON(s string) *JSON { return &JSON{json.RawMessage(s)} }

type Strings struct {
	a []string
}

func NewStrings() *Strings { return &Strings{[]string{}} }

func (ss *Strings) Append(val string) { ss.a = append(ss.a, val) }
func (ss *Strings) Get(i int) string  { return ss.a[i] }
func (ss *Strings) Length() int       { return len(ss.a) }
func (ss *Strings) String() string    { return fmt.Sprint(ss.a) }

type Props struct {
	m map[string]interface{}
}

func NewProps() *Props { return &Props{make(map[string]interface{})} }

func (ps *Props) String() string { return fmt.Sprint(ps.m) }

func (ps *Props) SetBool(key string, val bool)            { ps.m[key] = val }
func (ps *Props) SetInt(key string, val int)              { ps.m[key] = val }
func (ps *Props) SetFloat(key string, val float64)        { ps.m[key] = val }
func (ps *Props) SetString(key string, val string)        { ps.m[key] = val }
func (ps *Props) SetStringArray(key string, ref *Strings) { ps.m[key] = ref.a }
func (ps *Props) SetObject(key string, ref *Props)        { ps.m[key] = ref.m }
func (ps *Props) SetJSON(key string, ref *JSON)           { ps.m[key] = ref.x }

func (ps *Props) GetBool(key string) (val bool, err error) {
	if x, found := ps.m[key]; found {
		if b, ok := x.(bool); ok {
			val = b
		} else {
			err = fmt.Errorf("Prop type: %q is not a bool", key)
		}
	}
	return
}

func (ps *Props) GetInt(key string) (val int, err error) {
	if x, found := ps.m[key]; found {
		if f, ok := x.(float64); ok {
			val = int(f)
		} else {
			err = fmt.Errorf("Prop type: %q is not a number", key)
		}
	}
	return
}

func (ps *Props) GetFloat(key string) (val float64, err error) {
	if x, found := ps.m[key]; found {
		if f, ok := x.(float64); ok {
			val = f
		} else {
			err = fmt.Errorf("Prop type: %q is not a number", key)
		}
	}
	return
}

func (ps *Props) GetString(key string) (val string, err error) {
	if x, found := ps.m[key]; found {
		if s, ok := x.(string); ok {
			val = s
		} else {
			err = fmt.Errorf("Prop type: %q is not a string", key)
		}
	}
	return
}

func (ps *Props) GetStringArray(key string) (ref *Strings, err error) {
	if x, found := ps.m[key]; found {
		if xs, ok := x.([]interface{}); ok {
			ref = &Strings{make([]string, len(xs))}
			for i, x := range xs {
				if s, ok := x.(string); ok {
					ref.a[i] = s
				} else {
					err = fmt.Errorf("Prop type: %q is not a string array", key)
					return
				}
			}
		} else {
			err = fmt.Errorf("Prop type: %q is not an array", key)
		}
	}
	return
}

func (ps *Props) GetObject(key string) (ref *Props, err error) {
	if x, found := ps.m[key]; found {
		if m, ok := x.(map[string]interface{}); ok {
			ref = &Props{m}
		} else {
			err = fmt.Errorf("Prop type: %q is not an object", key)
		}
	}
	return
}

func (ps *Props) GetObjectArray(key string) (ref *Objects, err error) {
	if x, found := ps.m[key]; found {
		if xs, ok := x.([]interface{}); ok {
			ref = &Objects{make([]*Props, len(xs))}
			for i, x := range xs {
				if m, ok := x.(map[string]interface{}); ok {
					ref.a[i] = &Props{m}
				} else {
					err = fmt.Errorf("Prop type: %q is not an object array", key)
					return
				}
			}
		} else {
			err = fmt.Errorf("Prop type: %q is not an array", key)
		}
	}
	return
}

type Objects struct {
	a []*Props
}

func (os *Objects) Get(i int) *Props { return os.a[i] }
func (os *Objects) Length() int      { return len(os.a) }
func (os *Objects) String() string   { return fmt.Sprint(os.a) }

type PropVisitor interface {
	VisitBool(string, bool) error
	VisitNumber(string, float64) error
	VisitString(string, string) error
	VisitStringArray(string, *Strings) error
	VisitObject(string, *Props) error
	VisitObjectArray(string, *Objects) error
}

func (ps *Props) Accept(callback PropVisitor) (err error) {
	var (
		strings *Strings
		object  *Props
		objects *Objects
	)

	for k, x := range ps.m {
		switch v := x.(type) {
		case bool:
			err = callback.VisitBool(k, v)

		case float64:
			err = callback.VisitNumber(k, v)

		case string:
			err = callback.VisitString(k, v)

		case []interface{}:
			strings, err = ps.GetStringArray(k)
			if err == nil {
				err = callback.VisitStringArray(k, strings)
			} else {
				objects, err = ps.GetObjectArray(k)
				if err == nil {
					err = callback.VisitObjectArray(k, objects)
				}
			}

		case map[string]interface{}:
			object, err = ps.GetObject(k)
			if err == nil {
				err = callback.VisitObject(k, object)
			}
		}

		if err != nil {
			break
		}
	}

	return
}

type Payload struct {
	a []ninchat.Frame
}

func NewPayload() *Payload { return new(Payload) }

func (p *Payload) Append(blob []byte) {
	p.a = append(p.a, append([]byte{}, blob...))
}
func (p *Payload) Get(i int) []byte { return p.a[i] }
func (p *Payload) Length() int      { return len(p.a) }
func (p *Payload) String() string   { return fmt.Sprint(p.a) }

type SessionEventHandler interface {
	OnSessionEvent(params *Props)
}

type EventHandler interface {
	OnEvent(params *Props, payload *Payload, lastReply bool)
}

type CloseHandler interface {
	OnClose()
}

type ConnStateHandler interface {
	OnConnState(state string)
}

type ConnActiveHandler interface {
	OnConnActive()
}

type LogHandler interface {
	OnLog(msg string)
}

type session interface {
	SetParams(params map[string]interface{})
	Open()
	Close()
	Send(*ninchat.Action) error
}

type Session struct {
	session ninchat.Session
	s       session

	sessionEventHandler SessionEventHandler
	eventHandler        EventHandler
	closeHandler        CloseHandler
	connStateHandler    ConnStateHandler
	connActiveHandler   ConnActiveHandler
	logHandler          LogHandler

	closeCalled bool
}

func NewSession() (s *Session) {
	s = new(Session)

	s.session = ninchat.Session{
		OnSessionEvent: func(e *ninchat.Event) {
			s.sessionEventHandler.OnSessionEvent(&Props{e.Params})

			if s.closeCalled {
				s.releaseHandlers()
			}
		},

		OnEvent: func(e *ninchat.Event) {
			s.eventHandler.OnEvent(&Props{e.Params}, &Payload{e.Payload}, e.LastReply)
		},

		OnClose: func() {
			if s.closeHandler != nil {
				s.closeHandler.OnClose()
			}

			s.releaseHandlers()
		},

		OnConnState: func(state string) {
			if s.connStateHandler != nil {
				s.connStateHandler.OnConnState(state)
			}
		},

		OnConnActive: func() {
			if s.connActiveHandler != nil {
				s.connActiveHandler.OnConnActive()
			}
		},

		OnLog: func(fragments ...interface{}) {
			if s.logHandler != nil {
				var msg bytes.Buffer
				for i, x := range fragments {
					fmt.Fprint(&msg, x)
					if i < len(fragments)-1 {
						msg.WriteString(" ")
					}
				}
				s.logHandler.OnLog(msg.String())
			}
		},
	}

	s.s = ninchatstate.New(&s.session)
	return
}

// releaseHandlers attempts to free resources quickly in the iOS world.
func (s *Session) releaseHandlers() {
	s.sessionEventHandler = nil
	s.eventHandler = nil
	s.closeHandler = nil
	s.connStateHandler = nil
	s.connActiveHandler = nil
	s.logHandler = nil

	runtime.GC()
}

func (s *Session) SetOnSessionEvent(h SessionEventHandler) { s.sessionEventHandler = h }
func (s *Session) SetOnEvent(h EventHandler)               { s.eventHandler = h }
func (s *Session) SetOnClose(h CloseHandler)               { s.closeHandler = h }
func (s *Session) SetOnConnState(h ConnStateHandler)       { s.connStateHandler = h }
func (s *Session) SetOnConnActive(h ConnActiveHandler)     { s.connActiveHandler = h }
func (s *Session) SetOnLog(h LogHandler)                   { s.logHandler = h }

func (s *Session) SetAddress(address string) {
	s.session.Address = address
}

func (s *Session) SetParams(params *Props) (err error) {
	defer func() {
		err = asError(recover())
	}()

	s.s.SetParams(params.m)
	return
}

func (s *Session) Open() (err error) {
	defer func() {
		err = asError(recover())
	}()

	s.s.Open()
	return
}

func (s *Session) Close() {
	s.closeCalled = true
	s.s.Close()
}

func (s *Session) Send(params *Props, payload *Payload) (actionId int64, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = asError(x)
		}
	}()

	// This gotcha cannot be checked in the common Send function, in order to
	// retain bug-by-bug backward-compatibility.  But it can be enforced by
	// this new API.
	if x, found := params.m["action_id"]; found && x != nil {
		err = errors.New("action_id specified and is not null")
		return
	}

	action := &ninchat.Action{
		Params: params.m,
	}
	if payload != nil {
		action.Payload = payload.a
	}
	err = s.s.Send(action)
	if err == nil {
		actionId, _ = action.GetID()
	}
	return
}

type Event ninchat.Event

func (e *Event) GetParams() *Props    { return &Props{e.Params} }
func (e *Event) GetPayload() *Payload { return &Payload{e.Payload} }
func (e *Event) String() string       { return fmt.Sprint(*e) }

type Events struct {
	a []*ninchat.Event
}

func (es *Events) Get(i int) *Event { return (*Event)(es.a[i]) }
func (es *Events) Length() int      { return len(es.a) }
func (es *Events) String() string   { return fmt.Sprint(es.a) }

type Caller struct {
	c ninchat.Caller
}

func NewCaller() *Caller {
	return new(Caller)
}

func (c *Caller) SetAddress(address string) {
	c.c.Address = address
}

func (c *Caller) Call(params *Props, payload *Payload) (events *Events, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = asError(x)
		}
	}()

	action := &ninchat.Action{
		Params: params.m,
	}
	if payload != nil {
		action.Payload = payload.a
	}
	es, err := c.c.Call(action)
	if err == nil {
		events = &Events{es}
	}
	return
}
