// +build !js

package ninchat

import (
	"encoding/base64"
)

// longPollBinaryPayload converts binary data to base64 strings.
func longPollBinaryPayload(action *Action) (payload []interface{}) {
	for _, data := range action.Payload {
		text := base64.StdEncoding.EncodeToString(data)
		payload = append(payload, text)
	}
	return
}
