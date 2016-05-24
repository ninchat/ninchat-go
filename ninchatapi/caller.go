package ninchatapi

import (
	"github.com/ninchat/ninchat-go"
)

// Caller sends actions without a session.
type Caller struct {
	Params map[string]interface{} // Additional (authentication) parameters.
}

func (c *Caller) Send(action *ninchat.Action) (err error) {
	for name, value := range c.Params {
		if _, found := action.Params[name]; !found {
			action.Params[name] = value
		}
	}

	_, err = ninchat.Call(action)
	return
}

var (
	defaultCaller Caller
)
