// +build js

package ninchat

import (
	"github.com/gopherjs/gopherjs/js"
)

type Action struct {
	Params  map[string]interface{}
	Payload []*js.Object
	OnReply func(*Event)

	id int64
}

type Event struct {
	Params    map[string]interface{}
	Payload   []*js.Object
	LastReply bool
}
