package ninchatmessage

import (
	"encoding/json"
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

func (m *Link) Marshal() (payload [][]byte, err error) {
	return marshalJSON(m)
}

func (m *Link) Unmarshal(payload [][]byte) error {
	return json.Unmarshal(payload[0], m)
}
