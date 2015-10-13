// +build js

package ninchat

import (
	"github.com/gopherjs/gopherjs/js"
)

// longPollBinaryPayload converts data URIs to base64 strings.
func longPollBinaryPayload(action *Action) (payload []interface{}) {
	for _, dataURI := range action.Payload {
		base64 := (*js.Object)(dataURI).Call("split", ",").Index(1)
		payload = append(payload, base64)
	}
	return
}
