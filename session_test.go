package ninchat_test

import (
	"testing"

	"."
)

func TestSession(t *testing.T) {
	testSession(t, "")
}

func TestSessionLongpoll(t *testing.T) {
	testSession(t, "longpoll")
}

func testSession(t *testing.T, transport string) {
	events := make(chan *ninchat.Event, 10)

	session := &ninchat.Session{
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

	session.SetParams(map[string]interface{}{
		"message_types": []string{
			"*",
		},
	})

	session.SetTransport(transport)

	session.Open()
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

			session.Send(&ninchat.Action{
				Params: map[string]interface{}{
					"action":       "send_message",
					"user_id":      userId,
					"message_type": "ninchat.com/text",
				},
				Payload: [][]byte{
					[]byte("{\"text\":\"hello\"}"),
				},
			})

		case "message_received":
			session.Send(&ninchat.Action{
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
