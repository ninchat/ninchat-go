package webhook

import (
	"encoding/json"
)

type Webhook struct {
	// Kid and Exp are specified only interfaces which use signatures
	// (i.e. they are not specified for AWS Lambda functions).
	Kid string `json:"kid,omitempty"`
	Exp int64  `json:"exp,omitempty"`

	Aud     string    `json:"aud"`
	Event   EventType `json:"event"`
	EventID string    `json:"event_id"`

	// Raw JSON event fields can be used to store the data in lossless form.
	WebhookVerificationJSON json.RawMessage `json:"webhook_verification,omitempty"`
	AudienceAcceptedJSON    json.RawMessage `json:"audience_accepted,omitempty"`
	AudienceCompleteJSON    json.RawMessage `json:"audience_complete,omitempty"`
	DataAccessJSON          json.RawMessage `json:"data_access,omitempty"`
}

func (doc *Webhook) WebhookVerificationResponse() (content []byte) {
	content, err := json.Marshal(WebhookVerificationResponse{AudNinchat, doc.WebhookVerificationJSON})
	if err != nil {
		panic(err)
	}
	return
}

func (doc *Webhook) AudienceAccepted() (event AudienceAccepted, err error) {
	err = json.Unmarshal(doc.AudienceAcceptedJSON, &event)
	return
}

func (doc *Webhook) AudienceComplete() (event AudienceComplete, err error) {
	err = json.Unmarshal(doc.AudienceCompleteJSON, &event)
	return
}

func (doc *Webhook) DataAccess() (event DataAccess, err error) {
	err = json.Unmarshal(doc.DataAccessJSON, &event)
	return
}
