// +build !js

package ninchat

// Action
type Action struct {
	Params  map[string]interface{}
	Payload []Frame
	OnReply func(*Event)

	id int64
}

// Event
type Event struct {
	Params    map[string]interface{}
	Payload   []Frame
	LastReply bool
}

// Frame
type Frame []byte
