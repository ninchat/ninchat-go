// +build test

package ninchat

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"testing"
)

const numDeliverMessages = 10

var deliverMessageType = "test" + strconv.Itoa(rand.Int())

func init() {
	log.SetFlags(log.Flags() | log.Lmicroseconds)
}

func TestDialogueMessageDelivery(t *testing.T) {
	testMessageDeliveries(t, "")
}

func TestChannelMessageDelivery(t *testing.T) {
	id := os.Getenv("NINCHAT_TEST_CHANNEL_ID")
	if id == "" {
		t.Skip("NINCHAT_TEST_CHANNEL_ID environment variable not set")
	}

	testMessageDeliveries(t, id)
}

func testMessageDeliveries(t *testing.T, channelID string) {
	t.Run("Normal", func(t *testing.T) {
		testMessageDelivery(t, channelID, func(s *Session) func() {
			return func() {}
		})
	})

	t.Run("BreakConnOnly", func(t *testing.T) {
		testMessageDelivery(t, channelID, func(s *Session) func() {
			return s.TestBreakConnection()
		})
	})

	t.Run("BreakConnLoseSession", func(t *testing.T) {
		testMessageDelivery(t, channelID, func(s *Session) func() {
			s.TestLoseSession()
			return s.TestBreakConnection()
		})
	})

	t.Run("ForgetConnOnly", func(t *testing.T) {
		testMessageDelivery(t, channelID, func(s *Session) func() {
			return s.TestForgetConnection()
		})
	})

	t.Run("ForgetConnLoseSession", func(t *testing.T) {
		testMessageDelivery(t, channelID, func(s *Session) func() {
			s.TestLoseSession()
			return s.TestForgetConnection()
		})
	})
}

func testMessageDelivery(t *testing.T, channelID string, disconnect func(*Session) func()) {
	session1, events1 := openSessionForMessageDelivery(t, 1)
	defer session1.Close()
	uid1, _ := (<-events1).Str("user_id")

	session2, events2 := openSessionForMessageDelivery(t, 2)
	defer session2.Close()
	uid2, _ := (<-events2).Str("user_id")

	if channelID != "" {
		joinChannel(session1, events1, channelID, uid1)
		joinChannel(session2, events2, channelID, uid2)
	}

	seen2 := make(map[string]struct{})

	sendMessage(session1, channelID, uid2, 255)
	if val := receiveMessage(t, session2, events2, seen2, channelID, uid1); val != 255 {
		t.Fatal(val)
	}

	reconnect := disconnect(session2)
	if reconnect == nil {
		t.Fatal("disconnection failed")
	}

	for i := 0; i < numDeliverMessages; i++ {
		sendMessage(session1, channelID, uid2, i)
	}

	reconnect()

	var values []int
	for i := 0; i < numDeliverMessages; i++ {
		values = append(values, receiveMessage(t, session2, events2, seen2, channelID, uid1))
	}
	sort.Ints(values)
	for i, val := range values {
		if val != i {
			t.Fatal(values)
		}
	}

	sendMessage(session1, channelID, uid2, 254)
	if val := receiveMessage(t, session2, events2, seen2, channelID, uid1); val != 254 {
		t.Fatal(val)
	}
}

func openSessionForMessageDelivery(t *testing.T, index int) (*Session, <-chan *Event) {
	events := make(chan *Event, 1000)

	session := &Session{
		OnSessionEvent: func(e *Event) {
			events <- e
		},

		OnEvent: func(e *Event) {
			log.Printf("%d EVENT: %v %v\n", index, e.Payload, e.Params)
			events <- e
		},

		OnConnState: func(s string) {
			log.Printf("%d CONN STATE: %s\n", index, s)
		},

		OnLog: func(x ...interface{}) {
			log.Println(append([]interface{}{fmt.Sprintf("%d LOG:", index)}, x...)...)
		},
	}

	session.SetParams(map[string]interface{}{"message_types": []string{deliverMessageType}})
	session.Open()

	return session, events
}

func joinChannel(session *Session, events <-chan *Event, channelID, uid string) {
	a := &Action{
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

func sendMessage(session *Session, channelID, destUID string, val int) {
	params := map[string]interface{}{
		"action":       "send_message",
		"message_type": deliverMessageType,
	}
	if channelID == "" {
		params["user_id"] = destUID
	} else {
		params["channel_id"] = channelID
	}

	session.Send(&Action{
		Params:  params,
		Payload: []Frame{Frame{byte(val)}},
	})
}

func receiveMessage(t *testing.T, session *Session, events <-chan *Event, seen map[string]struct{}, channelID, srcUID string) int {
	t.Helper()

	for e := range events {
		switch e.String() {
		case "session_created":
			params := map[string]interface{}{
				"action":         "load_history",
				"history_length": numDeliverMessages,
			}
			if channelID == "" {
				params["user_id"] = srcUID
			} else {
				params["channel_id"] = channelID
			}

			session.Send(&Action{Params: params})

		case "message_received":
			if x, _ := e.Str("message_user_id"); x == srcUID {
				messageID, _ := e.Str("message_id")

				if _, found := seen[messageID]; !found {
					seen[messageID] = struct{}{}

					if x, _ := e.Str("message_type"); x != deliverMessageType {
						t.Fatal(x)
					}

					return int(e.Payload[0][0])
				}
			}
		}
	}

	panic("out of events")
}
