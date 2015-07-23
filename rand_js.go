// +build js

package ninchat

import (
	"github.com/gopherjs/gopherjs/js"
)

func randFloat64() float64 {
	return js.Global.Get("Math").Call("random").Float()
}
