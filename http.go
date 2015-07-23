package ninchat

func newJSONRequest(url string, action map[string]interface{}) (req *httpRequest, err error) {
	data, err := jsonMarshal(action)
	if err != nil {
		return
	}

	return newDataRequest("POST", url, data)
}

func getJSONRequestResponseChannel(url string, action map[string]interface{}, timeout duration) <-chan httpResponse {
	c := make(chan httpResponse, 1)

	if req, err := newJSONRequest(url, action); err == nil {
		putResponseToChannel(req, timeout, c)
	} else {
		c <- httpResponse{err: err}
	}

	return c
}

func getResponseChannel(req *httpRequest, timeout duration) <-chan httpResponse {
	c := make(chan httpResponse, 1)
	putResponseToChannel(req, timeout, c)
	return c
}
