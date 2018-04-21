package client

import (
	"bytes"
	"encoding/json"
	"fmt"

	ninchat "github.com/ninchat/ninchat-go"
)

type SessionEventJSONHandler interface {
	HandleSessionEventJSON(encodedParams []byte)
}

type EventJSON2Handler interface {
	HandleEventJSON2(encodedParams, payloadFrame1, payloadFrame2 []byte, lastReply bool)
}

type CloseHandler interface {
	HandleClose()
}

type ConnStateHandler interface {
	HandleConnState(state string)
}

type ConnActiveHandler interface {
	HandleConnActive()
}

type LogHandler interface {
	HandleLog(msg string)
}

type Session struct {
	ninchat.Session
}

func NewSession() *Session {
	return new(Session)
}

func (s *Session) OnSessionEventJSON(callback SessionEventJSONHandler) {
	s.Session.OnSessionEvent = func(e *ninchat.Event) {
		params, err := json.Marshal(e.Params)
		if err != nil {
			panic(err)
		}

		callback.HandleSessionEventJSON(params)
	}
}

func (s *Session) OnEventJSON2(callback EventJSON2Handler) {
	s.Session.OnEvent = func(e *ninchat.Event) {
		params, err := json.Marshal(e.Params)
		if err != nil {
			panic(err)
		}

		var (
			frame1 []byte
			frame2 []byte
		)
		if len(e.Payload) > 0 {
			frame1 = e.Payload[0]
			if len(e.Payload) > 1 {
				frame2 = e.Payload[1]
			}
		}

		callback.HandleEventJSON2(params, frame1, frame2, e.LastReply)
	}
}

func (s *Session) OnClose(callback CloseHandler) {
	s.Session.OnClose = callback.HandleClose
}

func (s *Session) OnConnState(callback ConnStateHandler) {
	s.Session.OnConnState = callback.HandleConnState
}

func (s *Session) OnConnActive(callback ConnActiveHandler) {
	s.Session.OnConnActive = callback.HandleConnActive
}

func (s *Session) OnLog(callback LogHandler) {
	s.Session.OnLog = func(fragments ...interface{}) {
		var msg bytes.Buffer

		for i, x := range fragments {
			fmt.Fprint(&msg, x)
			if i < len(fragments)-1 {
				msg.WriteString(" ")
			}
		}

		callback.HandleLog(msg.String())
	}
}

func (s *Session) SetAddress(address string) {
	s.Address = address
}

func (s *Session) SetParamsJSON(encodedParams []byte) (err error) {
	params := make(map[string]interface{})

	err = json.Unmarshal(encodedParams, params)
	if err != nil {
		return
	}

	s.SetParams(params)
	return
}

func (s *Session) SendJSON2(encodedParams, payloadFrame1, payloadFrame2 []byte) (err error) {
	params := make(map[string]interface{})

	err = json.Unmarshal(encodedParams, params)
	if err != nil {
		return
	}

	action := &ninchat.Action{
		Params: params,
	}

	if payloadFrame2 != nil {
		action.Payload = []ninchat.Frame{payloadFrame1, payloadFrame2}
	} else if payloadFrame1 != nil {
		action.Payload = []ninchat.Frame{payloadFrame1}
	}

	return s.Send(action)
}
