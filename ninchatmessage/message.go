// Standard Ninchat message data structures.
//
// https://ninchat.com/api/v2#message-types
//
package ninchatmessage

import (
	"encoding/json"
)

// Message is implemented by all message types in this package.
type Message interface {
	MessageType() string
	Marshal() (payload [][]byte, err error)
	Unmarshal(payload [][]byte) error
}

func marshalJSON(m Message) (payload [][]byte, err error) {
	data, err := json.Marshal(m)
	if err == nil {
		payload = [][]byte{data}
	}
	return
}
