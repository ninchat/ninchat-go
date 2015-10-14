// +build js

package ninchatmessage

import (
	"encoding/json"

	"github.com/gopherjs/gopherjs/js"
	"github.com/ninchat/ninchat-go"
)

func marshalJSON(c Content) (payload []ninchat.Frame, err error) {
	data, err := json.Marshal(c)
	if err != nil {
		return
	}

	x := js.Global.Get("Uint8Array").New(data)
	payload = []ninchat.Frame{x}
	return
}

func unmarshalJSON(payload []ninchat.Frame, c Content) error {
	jsString := ninchat.StringifyFrame(payload[0])
	data := []byte(jsString.String())
	return json.Unmarshal(data, c)
}
