// +build js

package ninchat

import (
	"errors"
)

func jsError(x interface{}) (err error) {
	if x != nil {
		switch t := x.(type) {
		case error:
			err = t

		case string:
			err = errors.New(t)

		default:
			err = errors.New("?")
		}
	}

	return
}
