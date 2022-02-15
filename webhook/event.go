package webhook

import (
	"encoding/json"
)

type EventType string

const (
	EventWebhookVerification    EventType = "webhook_verification"
	EventAudienceRequested      EventType = "audience_requested"
	EventAudienceRequestDropped EventType = "audience_request_dropped"
	EventAudienceAccepted       EventType = "audience_accepted"
	EventAudienceComplete       EventType = "audience_complete"
	EventMessageSent            EventType = "message_sent"
	EventDataAccess             EventType = "data_access"
)

type AudienceRequested struct {
	RealmID    string   `json:"realm_id"`
	QueueID    string   `json:"queue_id"`
	AudienceID string   `json:"audience_id"`
	Audience   Audience `json:"audience"`
}

type AudienceRequestDropped struct {
	RealmID    string   `json:"realm_id"`
	QueueID    string   `json:"queue_id"`
	AudienceID string   `json:"audience_id"`
	Audience   Audience `json:"audience"`
}

type AudienceAccepted struct {
	RealmID    string   `json:"realm_id"`
	QueueID    string   `json:"queue_id"`
	AudienceID string   `json:"audience_id"`
	Audience   Audience `json:"audience"`
	DialogueID []string `json:"dialogue_id,omitempty"`
	ChannelID  string   `json:"channel_id,omitempty"`
}

type AudienceComplete struct {
	RealmID               string              `json:"realm_id"`
	QueueID               string              `json:"queue_id"`
	AudienceID            string              `json:"audience_id"`
	Audience              Audience            `json:"audience"`
	DialogueID            []string            `json:"dialogue_id,omitempty"`
	ChannelID             string              `json:"channel_id,omitempty"`
	MemberMessageMetadata map[string]Metadata `json:"member_message_metadata"`
	Messages              []Message           `json:"messages,omitempty"`
}

type MessageSent struct {
	RealmID    string  `json:"realm_id"`
	QueueID    string  `json:"queue_id,omitempty"`
	AudienceID string  `json:"audience_id,omitempty"`
	ChannelID  string  `json:"channel_id,omitempty"`
	Message    Message `json:"audience"`
}

type DataAccess struct {
	RealmID   string `json:"realm_id,omitempty"`
	ChannelID string `json:"channel_id,omitempty"`
	Query     Query  `json:"query"`
}

type Audience struct {
	RequestTime  float64                   `json:"request_time,omitempty"` // Always present in event; never in response.
	DropTime     float64                   `json:"drop_time,omitempty"`
	AcceptTime   float64                   `json:"accept_time,omitempty"`
	FinishTime   float64                   `json:"finish_time,omitempty"`
	CompleteTime float64                   `json:"complete_time,omitempty"`
	Members      map[string]AudienceMember `json:"members,omitempty"` // Always present in event; never in response.
	Metadata     Metadata                  `json:"metadata,omitempty"`
}

type AudienceMember struct {
	Agent    bool `json:"agent,omitempty"`
	Customer bool `json:"customer,omitempty"`
}

type Message struct {
	ID       string      `json:"id"`
	Time     float64     `json:"time,omitempty"` // Always present in event, may be omitted in response.
	Type     MessageType `json:"type,omitempty"` // Always present in event, may be omitted in response.
	UserID   string      `json:"user_id,omitempty"`
	UserName *string     `json:"user_name,omitempty"`
	Fold     bool        `json:"fold,omitempty"`

	PayloadJSON json.RawMessage `json:"payload"`
}

type Query struct {
	AudienceMetadata      bool           `json:"audience.metadata,omitempty"`
	MemberMessageMetadata bool           `json:"member_message_metadata,omitempty"`
	Messages              *MessagesQuery `json:"messages,omitempty"`
}

type MessagesQuery struct {
	MinID string `json:"min_id,omitempty"`
	MaxID string `json:"max_id,omitempty"`
}
