package client

import (
	"bytes"
	"fmt"

	ninchat "github.com/ninchat/ninchat-go"
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

type Strings struct {
	a []string
}

func NewStrings() *Strings { return new(Strings) }

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

type Payload struct {
	a []ninchat.Frame
}

func NewPayload() *Payload { return new(Payload) }

func (p *Payload) Append(blob []byte) { p.a = append(p.a, blob) }
func (p *Payload) Get(i int) []byte   { return p.a[i] }
func (p *Payload) Length() int        { return len(p.a) }
func (p *Payload) String() string     { return fmt.Sprint(p.a) }

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

type Session struct {
	s ninchat.Session
}

func NewSession() *Session {
	return new(Session)
}

func (s *Session) SetOnSessionEvent(callback SessionEventHandler) {
	s.s.OnSessionEvent = func(e *ninchat.Event) {
		callback.OnSessionEvent(&Props{e.Params})
	}
}

func (s *Session) SetOnEvent(callback EventHandler) {
	s.s.OnEvent = func(e *ninchat.Event) {
		callback.OnEvent(&Props{e.Params}, &Payload{e.Payload}, e.LastReply)
	}
}

func (s *Session) SetOnClose(callback CloseHandler) {
	s.s.OnClose = callback.OnClose
}

func (s *Session) SetOnConnState(callback ConnStateHandler) {
	s.s.OnConnState = callback.OnConnState
}

func (s *Session) SetOnConnActive(callback ConnActiveHandler) {
	s.s.OnConnActive = callback.OnConnActive
}

func (s *Session) SetOnLog(callback LogHandler) {
	s.s.OnLog = func(fragments ...interface{}) {
		var msg bytes.Buffer
		for i, x := range fragments {
			fmt.Fprint(&msg, x)
			if i < len(fragments)-1 {
				msg.WriteString(" ")
			}
		}
		callback.OnLog(msg.String())
	}
}

func (s *Session) SetAddress(address string) {
	s.s.Address = address
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
	s.s.Close()
}

func (s *Session) Send(params *Props, payload *Payload) (err error) {
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
	err = s.s.Send(action)
	return
}
