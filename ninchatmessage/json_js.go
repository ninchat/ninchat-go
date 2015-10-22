// +build js

package ninchatmessage

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/ninchat/ninchat-go"
)

func marshal(obj map[string]interface{}) (payload []ninchat.Frame, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = x.(error) // TODO
		}
	}()

	data := js.Global.Get("JSON").Call("stringify", obj)
	payload = []ninchat.Frame{data}
	return
}

func unmarshal(payload []ninchat.Frame) (x interface{}, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = x.(error) // TODO
		}
	}()

	data := ninchat.StringifyFrame(payload[0])
	x = js.Global.Get("JSON").Call("parse", data).Interface()
	return
}
