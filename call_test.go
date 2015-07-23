package ninchat_test

import (
	"testing"

	"."
)

func TestCall(t *testing.T) {
	events, err := ninchat.Call(&ninchat.Action{
		Params: map[string]interface{}{
			"action": "create_user",
		},
		OnReply: func(*ninchat.Event) {
			t.Log("REPLY")
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, e := range events {
		t.Log("EVENT:", e.Params)
	}

	if len(events) == 0 {
		t.Fail()
	}
}
