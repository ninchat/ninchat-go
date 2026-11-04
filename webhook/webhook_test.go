package webhook_test

import (
	"encoding/json"
	"testing"

	"github.com/ninchat/ninchat-go/webhook"
)

func TestAudienceFinished(t *testing.T) {
	var document webhook.Webhook
	err := json.Unmarshal([]byte(`{
		"aud": "realm:example",
		"event": "audience_finished",
		"event_id": "example",
		"audience_finished": {
			"realm_id": "realm",
			"realm": {"name": "Example"},
			"queue_id": "queue",
			"queue": {"name": "Support"},
			"audience_id": "audience",
			"channel_id": "channel",
			"audience": {
				"request_time": 1700000000.25,
				"accept_time": 1700000010.5,
				"finish_time": 1700000100.75,
				"members": {"customer": {"customer": true}, "agent": {"agent": true}},
				"metadata": {"topic": "Example"}
			},
			"channel_members_metadata": {"agent": {"role": "support"}}
		}
	}`), &document)
	if err != nil {
		t.Fatal(err)
	}
	if document.Event != webhook.EventAudienceFinished {
		t.Fatal(document.Event)
	}

	event, err := document.AudienceFinished()
	if err != nil {
		t.Fatal(err)
	}
	if event.RealmID != "realm" || event.QueueID != "queue" || event.AudienceID != "audience" || event.ChannelID != "channel" {
		t.Error(event)
	}
	if event.Audience.RequestTime != 1700000000.25 || event.Audience.AcceptTime != 1700000010.5 || event.Audience.FinishTime != 1700000100.75 {
		t.Error(event.Audience)
	}
	if event.QueueMember != nil || event.Audience.CompleteTime != 0 {
		t.Error(event)
	}
	if !event.Audience.Members["customer"].Customer || !event.Audience.Members["agent"].Agent {
		t.Error(event.Audience.Members)
	}
	if string(event.Audience.Metadata["topic"]) != `"Example"` || string(event.ChannelMembersMetadata["agent"]["role"]) != `"support"` {
		t.Error(event)
	}
}

func TestAudienceFinishedError(t *testing.T) {
	for _, raw := range []string{"", "{", `{"audience":{"finish_time":"invalid"}}`} {
		document := webhook.Webhook{AudienceFinishedJSON: json.RawMessage(raw)}
		if _, err := document.AudienceFinished(); err == nil {
			t.Errorf("no error for %q", raw)
		}
	}
}
