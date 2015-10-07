package ninchatapi_test

import (
	"bytes"
	"testing"

	"."
	".."
)

var messageData = []byte("{\"text\":\"hello\"}")

var imageData = []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x01\x00\x00\x00\x01\x08\x00\x00\x00\x00:~\x9bU\x00\x00\x00\nIDAT\x08\xd7c`\x00\x00\x00\x02\x00\x01\xe2!\xbc3\x00\x00\x00\x00IEND\xaeB`\x82")

func openSession(t *testing.T, transport string, params map[string]interface{}) (session *ninchat.Session, events chan *ninchat.Event) {
	events = make(chan *ninchat.Event, 10)

	session = &ninchat.Session{
		OnSessionEvent: func(event *ninchat.Event) {
			t.Log("SESSION EVENT:", event.Params)
			events <- event
		},

		OnEvent: func(event *ninchat.Event) {
			t.Log("EVENT:", event.Params)
			events <- event
		},

		OnConnState: func(state string) {
			t.Log("CONN STATE:", state)
		},

		OnConnActive: func() {
			t.Log("CONN ACTIVE")
		},

		OnLog: func(fragments ...interface{}) {
			t.Log(append([]interface{}{"LOG:"}, fragments...)...)
		},
	}

	session.SetParams(params)
	session.SetTransport(transport)
	session.Open()
	return
}

func TestSession(t *testing.T) {
	testSession(t, "")
}

func TestSessionLongpoll(t *testing.T) {
	testSession(t, "longpoll")
}

func testSession(t *testing.T, transport string) {
	params := map[string]interface{}{
		"message_types": []string{
			"*",
		},
	}

	session, events := openSession(t, transport, params)
	defer session.Close()

	sessionEvent := <-events
	if sessionEvent.String() != "session_created" {
		t.Fail()
	}

	userId := sessionEvent.Params["user_id"].(string)
	userAuth := sessionEvent.Params["user_auth"].(string)

	messageType := "ninchat.com/text"

	messageEvent, err := (&ninchatapi.SendMessage{
		UserId:      &userId,
		MessageType: &messageType,
		Payload: [][]byte{
			messageData,
		},
	}).Call(session)
	if err != nil {
		t.Fatal(err)
	}

	if len(messageEvent.Payload) != 1 {
		t.Fatal("payload length")
	}

	if bytes.Compare(messageEvent.Payload[0], messageData) != 0 {
		t.Error("payload content")
	}

	userEvent, err := (&ninchatapi.UpdateUser{
		PayloadAttrs: []string{"icon"},
		Payload: [][]byte{
			imageData,
		},
	}).Call(session)
	if err != nil {
		t.Fatal(err)
	}

	if userEvent.UserAttrs == nil {
		t.Error("UserAttrs")
	} else if url := userEvent.UserAttrs.Iconurl; url == nil || *url == "" {
		t.Error("Iconurl")
	}

	if _, err = (&ninchatapi.DeleteUser{UserAuth: &userAuth}).Call(session); err != nil {
		t.Fatal(err)
	}
}
