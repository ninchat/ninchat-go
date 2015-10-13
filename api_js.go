// +build js

package ninchat

import (
	"github.com/gopherjs/gopherjs/js"
)

type Action struct {
	Params  map[string]interface{}
	Payload []Frame
	OnReply func(*Event)

	id int64
}

type Event struct {
	Params    map[string]interface{}
	Payload   []Frame
	LastReply bool
}

type Frame *js.Object

func singleFrame(x *js.Object) []Frame {
	return []Frame{x}
}
