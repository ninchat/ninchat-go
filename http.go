package ninchat

import (
	"fmt"
	"runtime"
)

var DefaultUserAgent = fmt.Sprintf("ninchat-go/1 (%s; %s)", runtime.GOOS, runtime.GOARCH)

func prepareHeader(custom map[string][]string) map[string][]string {
	h := make(map[string][]string)
	for k, vv := range custom {
		h[k] = vv
	}

	if _, found := h[userAgentHeader]; !found {
		h[userAgentHeader] = []string{DefaultUserAgent}
	}

	return h
}

func newJSONRequest(url string, header map[string][]string, action map[string]interface{}) (req *httpRequest, err error) {
	data, err := jsonMarshal(action)
	if err != nil {
		return
	}

	return newDataRequest("POST", url, header, data)
}

func getResponseChannel(req *httpRequest, timeout duration) <-chan httpResponse {
	c := make(chan httpResponse, 1)
	putResponseToChannel(req, timeout, c)
	return c
}
