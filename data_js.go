// +build js

package ninchat

import (
	"github.com/gopherjs/gopherjs/js"
)

func emptyData() *js.Object {
	return js.Global.Get("ArrayBuffer").New()
}

func dataLength(x *js.Object) int {
	return x.Length()
}

func stringData(x *js.Object) *js.Object {
	return StringifyFrame(x)
}

func dataString(x *js.Object) string {
	return x.String()
}

// StringifyFrame is only available with GopherJS.
func StringifyFrame(object *js.Object) *js.Object {
	if _, ok := object.Interface().(string); ok {
		return object
	}

	jsView := js.Global.Get("Uint8Array").New(object)
	goBytes := jsView.Interface().([]uint8)
	goString := string(goBytes)
	return js.Global.Get("String").New(goString)
}
