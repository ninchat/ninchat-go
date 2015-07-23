// +build !js

package ninchat

import (
	"encoding/json"
)

func jsonMarshal(v map[string]interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func jsonUnmarshalArray(data []byte, v *[]interface{}) (err error) {
	return json.Unmarshal(data, v)
}

func jsonUnmarshalObject(data []byte, v *map[string]interface{}) (err error) {
	return json.Unmarshal(data, v)
}
