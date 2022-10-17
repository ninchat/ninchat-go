// See https://github.com/ninchat/ninchat-go/blob/master/webhook/example/processor.go for an example.
package webhook

import (
	"encoding/json"
)

// WrappedWebhook is the alternative content format.
type WrappedWebhook struct {
	// If this is false, then the webhook request content was not wrapped.
	Wrapped bool `json:"wrapped"`

	Signature string `json:"signature"`
	JSON      []byte `json:"base64"`
}

// Webhook is the default content format.
type Webhook struct {
	// Kid and Exp are specified only with transports or formats which use
	// signatures (they are specified with HTTP, and with AWS Lambda if the
	// webhook was wrapped).
	Kid string `json:"kid,omitempty"`
	Exp int64  `json:"exp,omitempty"`

	Aud     string    `json:"aud"`
	Event   EventType `json:"event"`
	EventID string    `json:"event_id"`

	// Raw JSON event fields can be used to store the data in lossless form.
	WebhookVerificationJSON    json.RawMessage `json:"webhook_verification,omitempty"`
	AudienceRequestedJSON      json.RawMessage `json:"audience_requested,omitempty"`
	AudienceRequestDroppedJSON json.RawMessage `json:"audience_request_dropped,omitempty"`
	AudienceAcceptedJSON       json.RawMessage `json:"audience_accepted,omitempty"`
	AudienceAgentReleasedJSON  json.RawMessage `json:"audience_agent_released,omitempty"`
	AudienceCompleteJSON       json.RawMessage `json:"audience_complete,omitempty"`
	MessageSentJSON            json.RawMessage `json:"message_sent,omitempty"`
	DataAccessJSON             json.RawMessage `json:"data_access,omitempty"`
}

func (doc *Webhook) WebhookVerificationResponse() (content []byte) {
	content, err := json.Marshal(WebhookVerificationResponse{AudNinchat, doc.WebhookVerificationJSON})
	if err != nil {
		panic(err)
	}
	return
}

func (doc *Webhook) AudienceRequested() (event AudienceRequested, err error) {
	err = json.Unmarshal(doc.AudienceRequestedJSON, &event)
	return
}

func (doc *Webhook) AudienceRequestDropped() (event AudienceRequestDropped, err error) {
	err = json.Unmarshal(doc.AudienceRequestDroppedJSON, &event)
	return
}

func (doc *Webhook) AudienceAccepted() (event AudienceAccepted, err error) {
	err = json.Unmarshal(doc.AudienceAcceptedJSON, &event)
	return
}

func (doc *Webhook) AudienceAgentReleased() (event AudienceAgentReleased, err error) {
	err = json.Unmarshal(doc.AudienceAgentReleasedJSON, &event)
	return
}

func (doc *Webhook) AudienceComplete() (event AudienceComplete, err error) {
	err = json.Unmarshal(doc.AudienceCompleteJSON, &event)
	return
}

func (doc *Webhook) MessageSent() (event AudienceComplete, err error) {
	err = json.Unmarshal(doc.MessageSentJSON, &event)
	return
}

func (doc *Webhook) DataAccess() (event DataAccess, err error) {
	err = json.Unmarshal(doc.DataAccessJSON, &event)
	return
}
