// Standard Ninchat message data structures.
//
// https://ninchat.com/api/v2#message-types
//
package ninchatmessage

import (
	"encoding/json"

	"github.com/ninchat/ninchat-go"
)

// Content is implemented by all message types in this package.
type Content interface {
	MessageType() string
	Marshal() (payload []ninchat.Frame, err error)
	Unmarshal(payload []ninchat.Frame) error
}

func marshalJSON(c Content) (payload []ninchat.Frame, err error) {
	data, err := json.Marshal(c)
	if err != nil {
		return
	}

	payload = []ninchat.Frame{data}
	return
}

func unmarshalJSON(payload []ninchat.Frame, c Content) error {
	return json.Unmarshal(payload[0], c)
}
