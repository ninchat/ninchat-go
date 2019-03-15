package webhook

import (
	"encoding/json"
)

type EventType string

const (
	EventAudienceComplete EventType = "audience_complete"
)

type AudienceComplete struct {
	RealmID    string    `json:"realm_id"`
	QueueID    string    `json:"queue_id"`
	AudienceID string    `json:"audience_id"`
	Audience   Audience  `json:"audience"`
	DialogueID []string  `json:"dialogue_id,omitempty"`
	ChannelID  string    `json:"channel_id,omitempty"`
	Messages   []Message `json:"messages,omitempty"`
}

type Audience struct {
	RequestTime  float64                   `json:"request_time"`
	AcceptTime   float64                   `json:"accept_time,omitempty"`
	FinishTime   float64                   `json:"finish_time,omitempty"`
	CompleteTime float64                   `json:"complete_time"`
	Members      map[string]AudienceMember `json:"members"`
	Metadata     Metadata                  `json:"metadata,omitempty"`
}

type AudienceMember struct {
	Agent    bool `json:"agent,omitempty"`
	Customer bool `json:"customer,omitempty"`
}

type Message struct {
	ID       string      `json:"id"`
	Time     float64     `json:"time"`
	Type     MessageType `json:"type"`
	UserID   string      `json:"user_id,omitempty"`
	UserName *string     `json:"user_name,omitempty"`
	Fold     bool        `json:"fold,omitempty"`

	PayloadJSON json.RawMessage `json:"payload"`
}
