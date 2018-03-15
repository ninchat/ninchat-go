package ninchat

import (
	"bytes"
	"testing"
)

var messageData = []byte("{\"text\":\"hello\"}")

var imageData = []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x01\x00\x00\x00\x01\x08\x00\x00\x00\x00:~\x9bU\x00\x00\x00\nIDAT\x08\xd7c`\x00\x00\x00\x02\x00\x01\xe2!\xbc3\x00\x00\x00\x00IEND\xaeB`\x82")

func openSession(t *testing.T, forceLongPoll bool, params map[string]interface{}) (session *Session, events chan *Event) {
	events = make(chan *Event, 10)

	session = &Session{
		OnSessionEvent: func(event *Event) {
			t.Log("SESSION EVENT:", event.Params)
			events <- event
		},

		OnEvent: func(event *Event) {
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

		forceLongPoll: forceLongPoll,
	}

	session.SetParams(params)
	session.Open()
	return
}

func TestSession(t *testing.T) {
	testSession(t, false)
}

func TestSessionLongPoll(t *testing.T) {
	testSession(t, true)
}

func testSession(t *testing.T, forceLongPoll bool) {
	params := map[string]interface{}{
		"message_types": []string{
			"*",
		},
	}

	session, events := openSession(t, forceLongPoll, params)
	defer session.Close()

	var (
		userId   interface{}
		userAuth interface{}
	)

	for event := range events {
		switch event.String() {
		case "session_created":
			userId = event.Params["user_id"]
			userAuth = event.Params["user_auth"]

			session.Send(&Action{
				Params: map[string]interface{}{
					"action":       "send_message",
					"user_id":      userId,
					"message_type": "ninchat.com/text",
				},
				Payload: []Frame{
					messageData,
				},
			})

		case "message_received":
			if len(event.Payload) != 1 {
				t.Error("payload length")
			}

			if bytes.Compare(event.Payload[0], messageData) != 0 {
				t.Error("payload content")
			}

			session.Send(&Action{
				Params: map[string]interface{}{
					"action":        "update_user",
					"payload_attrs": []interface{}{"icon"},
				},
				Payload: []Frame{
					imageData,
				},
			})

		case "user_updated":
			if attrs, _ := event.Map("user_attrs"); attrs == nil {
				t.Error("user_attrs")
			} else if url, _ := attrs["iconurl"].(string); url == "" {
				t.Error("iconurl")
			}

			session.Send(&Action{
				Params: map[string]interface{}{
					"action":    "delete_user",
					"user_auth": userAuth,
				},
			})

		case "user_deleted":
			return
		}
	}
}

func TestLoginFailure(t *testing.T) {
	testLoginFailure(t, false)
}

func TestLoginFailureLongPoll(t *testing.T) {
	testLoginFailure(t, true)
}

func testLoginFailure(t *testing.T, forceLongPoll bool) {
	params := map[string]interface{}{
		"identity_type": "email",
		"identity_name": "nonexistent@example.invalid",
		"identity_auth": "asdf",
		"message_types": []string{},
	}

	session, events := openSession(t, forceLongPoll, params)
	defer session.Close()

	event := <-events
	if event.String() == "error" {
		if errorType, _ := event.Str("error_type"); errorType != "identity_not_found" {
			t.Error(errorType)
		}
	} else {
		t.Error(event.String())
	}
}
