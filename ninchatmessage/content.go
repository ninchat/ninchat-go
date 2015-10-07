// Standard Ninchat message data structures.
//
// https://ninchat.com/api/v2#message-types
//
package ninchatmessage

import (
	"encoding/json"
)

// Content is implemented by all message types in this package.
type Content interface {
	MessageType() string
	Marshal() (payload [][]byte, err error)
	Unmarshal(payload [][]byte) error
}

func marshalJSON(c Content) (payload [][]byte, err error) {
	data, err := json.Marshal(c)
	if err == nil {
		payload = [][]byte{data}
	}
	return
}
