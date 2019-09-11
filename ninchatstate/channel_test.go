// +build test

package ninchatstate

import (
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ninchat/ninchat-go"
)

func TestChannelState(t *testing.T) {
	channelID := os.Getenv("NINCHAT_TEST_CHANNEL_ID")
	if channelID == "" {
		t.Skip("NINCHAT_TEST_CHANNEL_ID environment variable not set")
	}

	session, events := openSession(t)
	defer session.Close()

	var timeout <-chan time.Time

	for {
		var open bool
		var e *ninchat.Event

		select {
		case e, open = <-events:
			if !open {
				return
			}

		case <-timeout:
			return
		}

		switch e.String() {
		case "session_created":
			if userChannels, found := e.Map("user_channels"); !found {
				t.Error("session_created event: user_channels parameter not found")
			} else if len(userChannels) != 0 {
				t.Error("session_created event: user_channels parameter not empty")
			}

			uid, _ := e.Str("user_id")
			joinChannel(session, events, channelID, uid)

			time.AfterFunc(time.Second*5, func() {
				session.TestLoseSession()
				reconnect := session.TestBreakConnection()
				time.Sleep(time.Second * 5)
				reconnect()
			})

		case "channel_joined", "channel_parted", "channel_updated", "message_received":
			data, _ := json.MarshalIndent(e.Params, "", "    ")
			log.Printf("EVENT: %s", data)
		}
	}
}

func openSession(t *testing.T) (*State, <-chan *ninchat.Event) {
	events := make(chan *ninchat.Event, 10)

	session := New(&ninchat.Session{
		OnSessionEvent: func(e *ninchat.Event) { events <- e },
		OnEvent:        func(e *ninchat.Event) { events <- e },

		OnConnState: func(s string) {
			log.Printf("CONN: %s", s)
		},

		OnLog: func(x ...interface{}) {
			log.Println(append([]interface{}{"LOG:"}, x...)...)
		},
	})

	params := map[string]interface{}{
		"message_types": []string{"test"},
	}

	session.SetParams(params)
	session.Open()

	return session, events
}

func joinChannel(session *State, events <-chan *ninchat.Event, channelID, uid string) {
	a := &ninchat.Action{
		Params: map[string]interface{}{
			"action":     "join_channel",
			"channel_id": channelID,
		},
	}

	session.Send(a)
	actionID, _ := a.GetID()

	for e := range events {
		if id, _ := e.Int64("action_id"); id == actionID {
			return
		}
	}

	panic("out of events")
}
