package ninchat

// parseDataURI strips data URI header.
func parseDataURI(uri string) (data string) {
	for i, x := range uri {
		if x == ',' {
			data = uri[i+1:]
			break
		}
	}
	return
}
