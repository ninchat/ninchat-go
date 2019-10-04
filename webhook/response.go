package webhook

import (
	"encoding/json"
)

const AudNinchat = "https://ninchat.com"

type WebhookVerificationResponse struct {
	Aud                     string          `json:"aud"`
	WebhookVerificationJSON json.RawMessage `json:"webhook_verification"`
}

type DataAccessResponse struct {
	Audience              *Audience           `json:"audience,omitempty"`
	MemberMessageMetadata map[string]Metadata `json:"member_message_metadata"`
	Messages              []Message           `json:"messages,omitempty"`
}
