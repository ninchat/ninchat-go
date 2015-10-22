// +build !js

package ninchatmessage

import (
	"encoding/json"

	"github.com/ninchat/ninchat-go"
)

func marshal(obj map[string]interface{}) (payload []ninchat.Frame, err error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return
	}

	payload = []ninchat.Frame{data}
	return
}

func unmarshal(payload []ninchat.Frame) (x interface{}, err error) {
	err = json.Unmarshal(payload[0], &x)
	return
}
