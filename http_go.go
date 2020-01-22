// +build !js

package ninchat

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"time"
)

const userAgentHeader = "User-Agent"

var tlsConfig = tls.Config{
	PreferServerCipherSuites: true,
}

func init() {
	t := http.DefaultTransport.(*http.Transport)
	if t.TLSClientConfig != nil {
		// In case a future Go version has it set.
		t.TLSClientConfig.PreferServerCipherSuites = tlsConfig.PreferServerCipherSuites
	} else {
		// Support testing on best effort basis.
		t.TLSClientConfig = &tlsConfig
	}
}

type httpRequest http.Request

func newGetRequest(url string, header map[string][]string) (*httpRequest, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = prepareHeader(header)
	return (*httpRequest)(req), err
}

func newDataRequest(method, url string, header map[string][]string, data []byte) (*httpRequest, error) {
	req, err := http.NewRequest(method, url, ioutil.NopCloser(bytes.NewReader(data)))
	if err != nil {
		return nil, err
	}

	req.Header = prepareHeader(header)
	return (*httpRequest)(req), err
}

type httpResponse struct {
	data []byte
	err  error
}

func getResponseData(req *httpRequest, timeout duration) (data []byte, err error) {
	client := http.Client{
		Timeout: time.Duration(timeout),
	}

	resp, err := client.Do((*http.Request)(req))
	if err != nil {
		return
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	return ioutil.ReadAll(resp.Body)
}

func putResponseToChannel(req *httpRequest, timeout duration, c chan<- httpResponse) {
	go func() {
		data, err := getResponseData(req, timeout)
		c <- httpResponse{
			data: data,
			err:  err,
		}
	}()
}
