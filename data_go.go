// +build !js

package ninchat

// emptyData constructs a zero-length payload.
func emptyData() []byte {
	return []byte{}
}

// singleData constructs a payload with one part.
func singleData(x []byte) [][]byte {
	return [][]byte{x}
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
