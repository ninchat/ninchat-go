// +build js

package ninchat

import (
	"errors"
)

func jsError(function string, x interface{}) (err error) {
	if x != nil {
		switch t := x.(type) {
		case error:
			err = errors.New(function + ": " + t.Error())

		case string:
			err = errors.New(function + ": " + t)

		default:
			err = errors.New(function + ": ?")
		}
	}

	return
}
