package main

import (
	"github.com/gopherjs/gopherjs/js"
)

// Log passes a log message to a callback, if a callback is defined.
func Log(logInvocationName string, onLog js.Object, tokens ...interface{}) {
	if onLog == nil || onLog.IsUndefined() || onLog.IsNull() {
		return
	}

	message := ""

	for _, x := range tokens {
		str := "?"

		if y, ok := x.(string); ok {
			str = y
		} else if y, ok := x.(error); ok {
			str = y.Error()
		}

		if len(message) > 0 {
			message += " "
		}

		message += str
	}

	for len(message) > 0 && message[len(message)-1] == ' ' {
		message = message[:len(message)-1]
	}

	jsInvoke(logInvocationName, onLog, message)
}
