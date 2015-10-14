// +build !js

package ninchatmessage

import (
	"encoding/json"

	"github.com/ninchat/ninchat-go"
)

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
