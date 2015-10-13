package ninchatmessage

import (
	"github.com/ninchat/ninchat-go"
)

const (
	LinkType = "ninchat.com/link"
)

// Link represents https://ninchat.com/link messages.
type Link struct {
	Name      string  `json:"name"`
	Size      int     `json:"size"`
	Icon      string  `json:"icon"`
	URL       string  `json:"url"`
	Thumbnail *string `json:"thumbnail,omitempty"`
}

func (*Link) MessageType() string {
	return LinkType
}

func (m *Link) Marshal() (payload []ninchat.Frame, err error) {
	return marshalJSON(m)
}

func (m *Link) Unmarshal(payload []ninchat.Frame) error {
	return unmarshalJSON(payload, m)
}
