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
	Unmarshal(payload []ninchat.Frame) error
}

// OriginalContent is implemented by non-system message types.
type OriginalContent interface {
	Content
	Marshal() (payload []ninchat.Frame, err error)
}

// UnaryContent is implemented by all message types in this package.
type UnaryContent interface {
	Content
	Init(payload interface{})
}
