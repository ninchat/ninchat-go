package ninchatmodel_test

import (
	"fmt"
	"testing"

	"github.com/tsavola/pointer"

	model "."
	"github.com/ninchat/ninchat-go"
	api "github.com/ninchat/ninchat-go/ninchatapi"
	"github.com/ninchat/ninchat-go/ninchatmessage"
)

func TestState(t *testing.T) {
	events, err := ninchat.Call(&ninchat.Action{
		Params: map[string]interface{}{
			"action": "create_user",
			"user_attrs": map[string]interface{}{
				"name": "Objective",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	targetUserId, _ := events[0].Str("user_id")

	done := make(chan struct{}, 2)

	state := model.State{
		Session: ninchat.Session{
			OnLog: func(fragments ...interface{}) {
				fmt.Println(append([]interface{}{"log:"}, fragments...)...)
			},
		},

		Self: model.SelfState{
			OnChange: func(change model.Change, user *model.User, auth string) {
				fmt.Println("self:", user.Id, auth)
			},
		},

		Users: model.UserState{
			OnChange: func(change model.Change, user *model.User) {
				fmt.Println("user:", user.Id, *user.Attrs.Name)
				done <- struct{}{}
			},
		},

		Settings: model.SettingsState{
			OnChange: func(change model.Change, settings map[string]interface{}) {
				fmt.Println("settings:", settings)
			},
		},

		Dialogues: model.DialogueState{
			OnChange: func(change model.Change, dialogue *model.Dialogue) {
				switch change {
				case model.Added:
					fmt.Println("dialogue-added:", dialogue.PeerId, "status", dialogue.Status)

				case model.Updated:
					fmt.Println("dialogue-updated:", dialogue.PeerId, "status", dialogue.Status)

				case model.Removed:
					fmt.Println("dialogue-removed:", dialogue.PeerId)
				}
			},

			Messages: model.MessageState{
				OnReceive: func(peerId string, e *api.MessageReceived) {
					fmt.Println("message:", e.MessageType, string(e.Payload()[0]))
					done <- struct{}{}
				},
			},
		},
	}

	state.SetParams(map[string]interface{}{
		"message_types": []interface{}{
			"*",
		},
	})

	if err := state.Open(); err != nil {
		t.Fatal(err)
	}
	defer state.Close()

	payload, err := (&ninchatmessage.Text{
		Text: "hello world",
	}).Marshal()
	if err != nil {
		t.Fatal(err)
	}

	if _, err := (&api.SendMessage{
		MessageType: pointer.String(ninchatmessage.TextType),
		UserId:      &targetUserId,
		Payload:     payload,
	}).Invoke(&state.Session); err != nil {
		t.Fatal(err)
	}

	<-done
	<-done
}
