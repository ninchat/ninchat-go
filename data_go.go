// +build !js

package ninchat

// emptyData constructs a zero-length payload.
func emptyData() []byte {
	return []byte{}
}

// dataLength returns the length of the payload.
func dataLength(x []byte) int {
	return len(x)
}

// stringData flattens the input.
func stringData(data []byte) []byte {
	return data
}

// dataString converts the input to a Go string.
func dataString(x []byte) string {
	return string(x)
}
