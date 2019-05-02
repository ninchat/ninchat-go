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
	AudienceAcceptedJSON json.RawMessage `json:"audience_accepted,omitempty"`
	AudienceCompleteJSON json.RawMessage `json:"audience_complete,omitempty"`
}

func (doc *Webhook) AudienceAccepted() (event AudienceAccepted, err error) {
	err = json.Unmarshal(doc.AudienceAcceptedJSON, &event)
	return
}

func (doc *Webhook) AudienceComplete() (event AudienceComplete, err error) {
	err = json.Unmarshal(doc.AudienceCompleteJSON, &event)
	return
}
