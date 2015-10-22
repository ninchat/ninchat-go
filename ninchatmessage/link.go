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
	Url       string  `json:"url"`
	Thumbnail *string `json:"thumbnail,omitempty"`
}

func (*Link) MessageType() string {
	return LinkType
}

func (m *Link) Marshal() (payload []ninchat.Frame, err error) {
	obj := map[string]interface{}{
		"name": m.Name,
		"size": m.Size,
		"icon": m.Icon,
		"url":  m.Url,
	}

	if m.Thumbnail != nil {
		obj["thumbnail"] = *m.Thumbnail
	}

	return marshal(obj)
}

func (m *Link) Unmarshal(payload []ninchat.Frame) (err error) {
	x, err := unmarshal(payload)
	if err == nil {
		m.Init(x)
	}
	return
}

func (m *Link) Init(payload interface{}) {
	if obj, ok := payload.(map[string]interface{}); ok {
		if x := obj["name"]; x != nil {
			m.Name, _ = x.(string)
		}

		if x := obj["size"]; x != nil {
			y, _ := x.(float64)
			m.Size = int(y)
		}

		if x := obj["icon"]; x != nil {
			m.Icon, _ = x.(string)
		}

		if x := obj["url"]; x != nil {
			m.Url, _ = x.(string)
		}

		if x := obj["thumbnail"]; x != nil {
			if y, ok := x.(string); ok {
				m.Thumbnail = &y
			}
		}
	}
}
