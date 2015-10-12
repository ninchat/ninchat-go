package ninchatmessage

import (
	"encoding/json"
)

const (
	TextType   = "ninchat.com/text"
	NoticeType = "ninchat.com/notice"
)

// Text represents https://ninchat.com/text messages.
type Text struct {
	Text string `json:"text"`
}

func (*Text) MessageType() string {
	return TextType
}

func (m *Text) Marshal() (payload [][]byte, err error) {
	return marshalJSON(m)
}

func (m *Text) Unmarshal(payload [][]byte) error {
	return json.Unmarshal(payload[0], m)
}

func (m *Text) String() string {
	return m.Text
}

// Notice represents https://ninchat.com/notice messages.  They are similar to
// https://ninchat.com/text messages.
type Notice Text

func (*Notice) MessageType() string {
	return NoticeType
}

func (m *Notice) Marshal() (payload [][]byte, err error) {
	return marshalJSON(m)
}

func (m *Notice) Unmarshal(payload [][]byte) error {
	return json.Unmarshal(payload[0], m)
}

func (m *Notice) String() string {
	return m.Text
}
