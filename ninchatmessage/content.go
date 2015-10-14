// Standard Ninchat message data structures.
//
// https://ninchat.com/api/v2#message-types
//
package ninchatmessage

import (
	"github.com/ninchat/ninchat-go"
)

// Content is implemented by all message types in this package.
type Content interface {
	MessageType() string
	Marshal() (payload []ninchat.Frame, err error)
	Unmarshal(payload []ninchat.Frame) error
}
