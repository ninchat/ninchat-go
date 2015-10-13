package ninchatmessage

import (
	"github.com/ninchat/ninchat-go"
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

func (m *Text) Marshal() (payload []ninchat.Frame, err error) {
	return marshalJSON(m)
}

func (m *Text) Unmarshal(payload []ninchat.Frame) error {
	return unmarshalJSON(payload, m)
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

func (m *Notice) Marshal() (payload []ninchat.Frame, err error) {
	return marshalJSON(m)
}

func (m *Notice) Unmarshal(payload []ninchat.Frame) error {
	return unmarshalJSON(payload, m)
}

func (m *Notice) String() string {
	return m.Text
}
