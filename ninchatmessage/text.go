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
	return marshal(map[string]interface{}{
		"text": m.Text,
	})
}

func (m *Text) Unmarshal(payload []ninchat.Frame) (err error) {
	x, err := unmarshal(payload)
	if err == nil {
		m.Init(x)
	}
	return
}

func (m *Text) Init(payload interface{}) {
	if obj, ok := payload.(map[string]interface{}); ok {
		if x := obj["text"]; x != nil {
			m.Text, _ = x.(string)
		}
	}
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

func (m *Notice) Marshal() ([]ninchat.Frame, error) {
	return (*Text)(m).Marshal()
}

func (m *Notice) Unmarshal(payload []ninchat.Frame) error {
	return (*Text)(m).Unmarshal(payload)
}

func (m *Notice) Init(payload interface{}) {
	(*Text)(m).Init(payload)
}

func (m *Notice) String() string {
	return (*Text)(m).String()
}
